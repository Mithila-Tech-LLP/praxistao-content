# Chapter 39: Mini Project 3 — Algorithm Playground

This mini-project brings together the algorithms from the past six chapters: sorting, searching, dynamic programming, greedy, and string algorithms. You'll build a command-line algorithm benchmark and comparison tool that runs algorithms on generated datasets, measures performance, and produces reports.

## What You'll Build

A CLI tool that:
1. Runs multiple algorithms on the same input and compares results
2. Benchmarks performance with realistic data sizes
3. Visualizes algorithm complexity through measured timing
4. Validates correctness automatically

**Project structure:**
```
algo-playground/
├── main.go
├── sort/
│   ├── sort.go          — sorting algorithm implementations
│   └── sort_test.go
├── search/
│   ├── search.go        — search algorithm implementations
│   └── search_test.go
├── dp/
│   ├── dp.go            — DP solutions
│   └── dp_test.go
├── bench/
│   └── bench.go         — benchmarking harness
├── gen/
│   └── gen.go           — test data generators
└── report/
    └── report.go        — timing and result display
```

---

## gen/gen.go — Data Generators

```go
package gen

import (
    "math/rand"
    "sort"
)

// Random generates a random integer slice.
func Random(n, max int, seed int64) []int {
    r := rand.New(rand.NewSource(seed))
    arr := make([]int, n)
    for i := range arr { arr[i] = r.Intn(max) }
    return arr
}

// Sorted generates a sorted slice.
func Sorted(n int) []int {
    arr := make([]int, n)
    for i := range arr { arr[i] = i }
    return arr
}

// ReverseSorted generates a reverse-sorted slice.
func ReverseSorted(n int) []int {
    arr := make([]int, n)
    for i := range arr { arr[i] = n - i }
    return arr
}

// NearlySorted generates a nearly sorted slice (5% elements out of place).
func NearlySorted(n int, seed int64) []int {
    arr := Sorted(n)
    r := rand.New(rand.NewSource(seed))
    swaps := n / 20
    for i := 0; i < swaps; i++ {
        a, b := r.Intn(n), r.Intn(n)
        arr[a], arr[b] = arr[b], arr[a]
    }
    return arr
}

// ManyDuplicates generates a slice where values repeat heavily.
func ManyDuplicates(n int, seed int64) []int {
    r := rand.New(rand.NewSource(seed))
    arr := make([]int, n)
    for i := range arr { arr[i] = r.Intn(n / 10) }
    return arr
}

// RandomStrings generates random strings of given length.
func RandomStrings(count, length int, seed int64) []string {
    r := rand.New(rand.NewSource(seed))
    const letters = "abcdefghijklmnopqrstuvwxyz"
    strs := make([]string, count)
    for i := range strs {
        b := make([]byte, length)
        for j := range b { b[j] = letters[r.Intn(len(letters))] }
        strs[i] = string(b)
    }
    return strs
}

// Copy returns a copy of the slice.
func Copy(arr []int) []int {
    cp := make([]int, len(arr))
    copy(cp, arr)
    return cp
}

// IsStrictlySorted verifies correct ascending sort.
func IsStrictlySorted(arr []int) bool {
    for i := 1; i < len(arr); i++ {
        if arr[i] < arr[i-1] { return false }
    }
    return true
}

// MatchesStdlib verifies that our sort matches sort.Ints.
func MatchesStdlib(original, result []int) bool {
    ref := Copy(original)
    sort.Ints(ref)
    if len(ref) != len(result) { return false }
    for i := range ref {
        if ref[i] != result[i] { return false }
    }
    return true
}
```

---

## sort/sort.go — Sorting Implementations

```go
package sortalgos

// All algorithms sort in-place unless noted.

func BubbleSort(arr []int) {
    n := len(arr)
    for i := 0; i < n-1; i++ {
        swapped := false
        for j := 0; j < n-1-i; j++ {
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j]
                swapped = true
            }
        }
        if !swapped { break }
    }
}

func InsertionSort(arr []int) {
    for i := 1; i < len(arr); i++ {
        key := arr[i]
        j := i - 1
        for j >= 0 && arr[j] > key { arr[j+1] = arr[j]; j-- }
        arr[j+1] = key
    }
}

func MergeSort(arr []int) {
    if len(arr) <= 1 { return }
    mid := len(arr) / 2
    MergeSort(arr[:mid])
    MergeSort(arr[mid:])
    mergeInPlace(arr, mid)
}

func mergeInPlace(arr []int, mid int) {
    tmp := make([]int, len(arr))
    copy(tmp, arr)
    i, j, k := 0, mid, 0
    for i < mid && j < len(arr) {
        if tmp[i] <= tmp[j] { arr[k] = tmp[i]; i++ } else { arr[k] = tmp[j]; j++ }
        k++
    }
    for i < mid { arr[k] = tmp[i]; i++; k++ }
    for j < len(arr) { arr[k] = tmp[j]; j++; k++ }
}

func QuickSort(arr []int) { quickSort(arr, 0, len(arr)-1) }

func quickSort(arr []int, lo, hi int) {
    if lo < hi {
        p := partition(arr, lo, hi)
        quickSort(arr, lo, p-1)
        quickSort(arr, p+1, hi)
    }
}

func partition(arr []int, lo, hi int) int {
    // Random pivot to avoid worst-case on sorted input:
    mid := lo + (hi-lo)/2
    arr[mid], arr[hi] = arr[hi], arr[mid]
    pivot := arr[hi]
    i := lo - 1
    for j := lo; j < hi; j++ {
        if arr[j] <= pivot { i++; arr[i], arr[j] = arr[j], arr[i] }
    }
    arr[i+1], arr[hi] = arr[hi], arr[i+1]
    return i + 1
}

func HeapSort(arr []int) {
    n := len(arr)
    for i := n/2 - 1; i >= 0; i-- { heapify(arr, n, i) }
    for end := n - 1; end > 0; end-- {
        arr[0], arr[end] = arr[end], arr[0]
        heapify(arr, end, 0)
    }
}

func heapify(arr []int, n, i int) {
    largest := i
    l, r := 2*i+1, 2*i+2
    if l < n && arr[l] > arr[largest] { largest = l }
    if r < n && arr[r] > arr[largest] { largest = r }
    if largest != i {
        arr[i], arr[largest] = arr[largest], arr[i]
        heapify(arr, n, largest)
    }
}

func CountingSort(arr []int, maxVal int) []int {
    count := make([]int, maxVal+1)
    for _, v := range arr { count[v]++ }
    for i := 1; i <= maxVal; i++ { count[i] += count[i-1] }
    output := make([]int, len(arr))
    for i := len(arr) - 1; i >= 0; i-- {
        output[count[arr[i]]-1] = arr[i]
        count[arr[i]]--
    }
    return output
}
```

---

## bench/bench.go — Benchmarking Harness

```go
package bench

import (
    "fmt"
    "time"
)

// Result holds a single benchmark result.
type Result struct {
    Name    string
    N       int
    Input   string
    Elapsed time.Duration
    Correct bool
}

// Runner executes a sort function and measures its time.
func Runner(
    name string,
    n int,
    inputDesc string,
    prepare func() []int,      // Returns a fresh copy of input
    run func([]int),            // Sorts in-place
    verify func([]int, []int) bool,  // verify(original, sorted)
) Result {
    original := prepare()
    input := make([]int, len(original))
    copy(input, original)

    start := time.Now()
    run(input)
    elapsed := time.Since(start)

    return Result{
        Name:    name,
        N:       n,
        Input:   inputDesc,
        Elapsed: elapsed,
        Correct: verify(original, input),
    }
}

// MultiRun runs an algorithm multiple times and returns the median duration.
func MultiRun(
    name string,
    n int,
    inputDesc string,
    prepare func() []int,
    run func([]int),
    verify func([]int, []int) bool,
    times int,
) Result {
    durations := make([]time.Duration, times)
    var correct bool

    for i := range durations {
        original := prepare()
        input := make([]int, len(original))
        copy(input, original)

        start := time.Now()
        run(input)
        durations[i] = time.Since(start)

        if i == 0 { correct = verify(original, input) }
    }

    // Median:
    median := medianDuration(durations)
    return Result{Name: name, N: n, Input: inputDesc, Elapsed: median, Correct: correct}
}

func medianDuration(d []time.Duration) time.Duration {
    // Simple insertion sort on small slice:
    for i := 1; i < len(d); i++ {
        key := d[i]
        j := i - 1
        for j >= 0 && d[j] > key { d[j+1] = d[j]; j-- }
        d[j+1] = key
    }
    return d[len(d)/2]
}
```

---

## report/report.go — Result Display

```go
package report

import (
    "fmt"
    "strings"
    "time"

    "algo-playground/bench"
)

// PrintTable prints a formatted comparison table.
func PrintTable(results []bench.Result) {
    headers := []string{"Algorithm", "N", "Input Type", "Time", "Status"}
    rows := make([][]string, len(results))

    for i, r := range results {
        status := "OK"
        if !r.Correct { status = "WRONG" }
        rows[i] = []string{
            r.Name,
            fmt.Sprintf("%d", r.N),
            r.Input,
            formatDuration(r.Elapsed),
            status,
        }
    }

    widths := make([]int, len(headers))
    for i, h := range headers { widths[i] = len(h) }
    for _, row := range rows {
        for i, cell := range row {
            if len(cell) > widths[i] { widths[i] = len(cell) }
        }
    }

    printRow := func(cells []string) {
        for i, cell := range cells {
            fmt.Printf("%-*s  ", widths[i], cell)
        }
        fmt.Println()
    }
    printSep := func() {
        for _, w := range widths { fmt.Print(strings.Repeat("-", w+2)) }
        fmt.Println()
    }

    printSep()
    printRow(headers)
    printSep()
    for _, row := range rows { printRow(row) }
    printSep()
}

func formatDuration(d time.Duration) string {
    switch {
    case d < time.Microsecond:
        return fmt.Sprintf("%dns", d.Nanoseconds())
    case d < time.Millisecond:
        return fmt.Sprintf("%.2fµs", float64(d.Nanoseconds())/1000)
    case d < time.Second:
        return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1e6)
    default:
        return fmt.Sprintf("%.3fs", d.Seconds())
    }
}

// PrintComplexityEstimate infers complexity from two timing results.
func PrintComplexityEstimate(small, large bench.Result) {
    ratio := float64(large.Elapsed) / float64(small.Elapsed)
    sizeRatio := float64(large.N) / float64(small.N)

    var estimate string
    switch {
    case ratio < sizeRatio*1.2:
        estimate = "O(n)"
    case ratio < sizeRatio*float64(log2(large.N))*1.2:
        estimate = "O(n log n)"
    case ratio < sizeRatio*sizeRatio*1.2:
        estimate = "O(n²)"
    default:
        estimate = "> O(n²)"
    }

    fmt.Printf("\n%s: %v (n=%d) → %v (n=%d), ratio=%.1fx → estimated %s\n",
        large.Name, small.Elapsed, small.N, large.Elapsed, large.N, ratio, estimate)
}

func log2(n int) float64 {
    result := 0.0
    for n > 1 { n /= 2; result++ }
    return result
}
```

---

## main.go — CLI Entry Point

```go
package main

import (
    "fmt"
    "os"
    "strconv"

    sortalgos "algo-playground/sort"
    "algo-playground/bench"
    "algo-playground/gen"
    "algo-playground/report"
)

func main() {
    n := 10_000
    if len(os.Args) > 1 {
        if v, err := strconv.Atoi(os.Args[1]); err == nil { n = v }
    }

    fmt.Printf("=== Sorting Benchmark (n=%d) ===\n\n", n)

    type sortFn struct {
        name string
        fn   func([]int)
    }

    algos := []sortFn{
        {"InsertionSort", sortalgos.InsertionSort},
        {"MergeSort", sortalgos.MergeSort},
        {"QuickSort", sortalgos.QuickSort},
        {"HeapSort", sortalgos.HeapSort},
        {"stdlib sort.Ints", func(arr []int) {
            import "sort"
            sort.Ints(arr)
        }},
    }

    inputTypes := []struct {
        name    string
        prepare func() []int
    }{
        {"random", func() []int { return gen.Random(n, n*10, 42) }},
        {"sorted", func() []int { return gen.Sorted(n) }},
        {"reverse", func() []int { return gen.ReverseSorted(n) }},
        {"nearly-sorted", func() []int { return gen.NearlySorted(n, 42) }},
        {"many-dupes", func() []int { return gen.ManyDuplicates(n, 42) }},
    }

    // Skip InsertionSort for large n (O(n²) would be too slow):
    skipLargeN := map[string]bool{"InsertionSort": n > 50_000}

    var results []bench.Result
    for _, input := range inputTypes {
        for _, algo := range algos {
            if skipLargeN[algo.name] { continue }
            result := bench.MultiRun(
                algo.name, n, input.name,
                input.prepare,
                algo.fn,
                func(orig, sorted []int) bool { return gen.MatchesStdlib(orig, sorted) },
                5,
            )
            results = append(results, result)
        }
    }

    report.PrintTable(results)
}
```

---

## sort/sort_test.go — Tests

```go
package sortalgos_test

import (
    "math/rand"
    "sort"
    "testing"

    sortalgos "algo-playground/sort"
    "algo-playground/gen"
)

type sortFunc struct {
    name string
    fn   func([]int)
}

var algorithms = []sortFunc{
    {"BubbleSort", sortalgos.BubbleSort},
    {"InsertionSort", sortalgos.InsertionSort},
    {"MergeSort", sortalgos.MergeSort},
    {"QuickSort", sortalgos.QuickSort},
    {"HeapSort", sortalgos.HeapSort},
}

func TestAllSorts(t *testing.T) {
    cases := []struct {
        name  string
        input []int
    }{
        {"empty", []int{}},
        {"single", []int{1}},
        {"two-sorted", []int{1, 2}},
        {"two-reversed", []int{2, 1}},
        {"duplicates", []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}},
        {"all-same", []int{7, 7, 7, 7, 7}},
        {"already-sorted", gen.Sorted(100)},
        {"reverse-sorted", gen.ReverseSorted(100)},
        {"random-100", gen.Random(100, 1000, 42)},
        {"random-1000", gen.Random(1000, 10000, 99)},
    }

    for _, tc := range cases {
        tc := tc
        for _, algo := range algorithms {
            algo := algo
            t.Run(algo.name+"/"+tc.name, func(t *testing.T) {
                t.Parallel()
                input := gen.Copy(tc.input)
                ref := gen.Copy(tc.input)
                sort.Ints(ref)

                algo.fn(input)

                for i := range input {
                    if input[i] != ref[i] {
                        t.Errorf("index %d: got %d, want %d (full: %v → %v, want %v)",
                            i, input[i], ref[i], tc.input, input, ref)
                        return
                    }
                }
            })
        }
    }
}

func BenchmarkSorts(b *testing.B) {
    sizes := []int{100, 1000, 10000}

    for _, n := range sizes {
        for _, algo := range algorithms {
            if algo.name == "BubbleSort" && n > 1000 { continue }

            b.Run(algo.name+"/random/"+strconv.Itoa(n), func(b *testing.B) {
                original := gen.Random(n, n*10, 42)
                b.ResetTimer()
                for i := 0; i < b.N; i++ {
                    input := gen.Copy(original)
                    algo.fn(input)
                }
            })
        }
    }
}
```

---

## Running the Project

```bash
# Initialize module:
mkdir algo-playground && cd algo-playground
go mod init algo-playground

# Create the package structure above and run:
go run . 10000

# Run tests with race detector:
go test ./... -race

# Run benchmarks:
go test ./sort/... -bench=. -benchmem -count=3

# Compare algorithms across sizes to observe complexity curves:
go run . 1000
go run . 10000
go run . 100000

# Observe that:
# - InsertionSort time grows quadratically (skip for n=100000)
# - MergeSort/QuickSort/HeapSort grow as n log n
# - InsertionSort beats QuickSort for n=100 due to cache effects
# - QuickSort beats MergeSort at large n due to better constants
```

---

## Extension Challenges

### Beginner
1. Add `SelectionSort` to the sort package and include it in the benchmark. Verify it produces the correct output but has no adaptive behavior (same time on sorted/random input).

### Intermediate
2. Add a **string sorting benchmark**: sort slices of random strings of varying lengths using `sort.Strings`, a radix sort variant, and a trie-based sort. At what string length/count does radix sort win?
3. Add a **DP benchmark**: compare `FibNaive`, `FibMemo`, and `FibOptimal` for n=10, 20, 30, 40. `FibNaive(40)` should take noticeably longer. Plot the timing ratios.

### Advanced
4. Add a **parallel merge sort** that uses goroutines to sort sub-arrays concurrently. Below a threshold, fall back to sequential sort. Benchmark with n=1M. Find the threshold where parallelism pays off. Is the speedup proportional to CPU count?
5. Implement an **adaptive sorter** that detects the input type (sorted, reversed, random, many duplicates) using a statistical test (sample 50 random pairs and measure sortedness), then selects the best algorithm for that distribution. Benchmark it against `sort.Ints` to see how close you can get.
