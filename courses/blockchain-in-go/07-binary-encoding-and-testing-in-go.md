# Chapter 07: Binary Encoding and Testing in Go

Blocks and transactions are Go structs while they live in memory, but the moment they need to be hashed, signed, sent over a network, or saved to disk, they have to become plain bytes. This chapter compares the three ways GoChain could do that, explains why the project starts with one specific choice, and sets up the testing discipline every GoChain package follows from this point forward.

## Table of Contents

1. [Why Everything Eventually Becomes Bytes](#1-why-everything-eventually-becomes-bytes)
2. [Option One: encoding/json](#2-option-one-encodingjson)
3. [Option Two: encoding/gob](#3-option-two-encodinggob)
4. [Option Three: A Hand-Rolled Binary Format](#4-option-three-a-hand-rolled-binary-format)
5. [Comparing the Three, and Why GoChain Starts With gob](#5-comparing-the-three-and-why-gochain-starts-with-gob)
6. [Go's Testing Framework: go test and testing.T](#6-gos-testing-framework-go-test-and-testingt)
7. [Table-Driven Tests](#7-table-driven-tests)
8. [A Full Worked Test Suite](#8-a-full-worked-test-suite)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why Everything Eventually Becomes Bytes

A Go struct like our `Envelope` from Chapter 04 lives happily in your computer's memory as a structured, typed value while your program runs. But several things GoChain must do simply cannot operate on a live, in-memory struct at all:

- **Hashing** (Volume 2) requires a specific, fixed sequence of bytes to feed into a hash function — `crypto/sha256.Sum256` takes a `[]byte`, not a struct.
- **Sending data over a network connection** (Volume 7) means writing raw bytes onto a TCP socket, and the receiving computer — quite possibly a different machine with different internal memory layout — has to reconstruct an equivalent struct purely from those bytes.
- **Saving data to disk** (Volume 8, and even the simple flat-file storage in Volume 3) means writing bytes to a file, since a file is fundamentally just a sequence of bytes.

**Serialization** is the general term for converting a structured, in-memory value into a flat sequence of bytes (or, in some formats, text) that can be stored or transmitted; **deserialization** is the reverse process, reconstructing the original structured value from those bytes. Go's standard library offers several ready-made serialization approaches, and this chapter compares three of them concretely, using GoChain's own future `core.Transaction` type as the running example.

One subtlety worth flagging immediately, because it matters enormously once hashing enters the picture in Volume 2: serialization must be **canonical** — the same logical value must always produce the *exact same* bytes, every time, on every machine. If two structurally identical transactions could serialize to two different byte sequences (say, because a map's keys were iterated in a different order on two different runs), they would hash differently, breaking the tamper-evidence property from Chapter 01 entirely. Chapter 09 returns to this point in depth once real hashing is introduced; for now, keep it in mind as you compare the three formats below.

---

## 2. Option One: encoding/json

`encoding/json` is Go's standard library package for converting Go values to and from JSON (JavaScript Object Notation), a widely used, human-readable text format. It requires no external dependency and is extremely easy to use:

```go
package main

import (
	"encoding/json"
	"fmt"
)

// toyTransaction is a simplified stand-in for core.Transaction, used
// purely to compare encoding formats before the real type exists.
type toyTransaction struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"` // amount in gochips
}

func main() {
	tx := toyTransaction{ID: "abc123", Amount: 50}

	encoded, err := json.Marshal(tx)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(encoded))
	// {"id":"abc123","amount":50}

	var decoded toyTransaction
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", decoded)
	// {ID:abc123 Amount:50}
}
```

`json.Marshal(tx)` serializes `tx` into a `[]byte` containing its JSON text representation, and `json.Unmarshal(encoded, &decoded)` does the reverse, filling in `decoded` from those bytes (note the `&decoded` — `Unmarshal` needs a pointer so it can write into the value, exactly the pointer reasoning from Chapter 04, Section 4). The backtick-quoted text after each field, like `` `json:"id"` ``, is a **struct tag** — metadata attached to a field that tells `encoding/json` exactly what key name to use in the JSON output, rather than defaulting to the Go field name itself.

JSON's biggest strength is **human readability** — you can open a JSON file in any text editor, or print it straight to a terminal, and immediately understand what it contains, which makes it excellent for debugging, for GoChain's future JSON-RPC and REST API (Volume 10), and for configuration files. Its weaknesses, for GoChain's core, performance-sensitive data path specifically: it's verbose (field names are repeated as text in every single encoded value, adding up quickly across millions of transactions), and encoding/decoding text involves real parsing overhead compared to a format that maps more directly onto raw bytes.

---

## 3. Option Two: encoding/gob

`encoding/gob` is a serialization format specific to Go itself, also part of the standard library, designed for Go programs talking to other Go programs. It isn't human-readable, and it isn't meant to be — its entire design goal is speed and simplicity when both ends of a serialization are Go code.

```go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type toyTransaction struct {
	ID     string
	Amount int64 // amount in gochips
}

func main() {
	tx := toyTransaction{ID: "abc123", Amount: 50}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(tx); err != nil {
		panic(err)
	}
	fmt.Printf("encoded gob bytes (not human-readable): %d bytes\n", buf.Len())

	var decoded toyTransaction
	decoder := gob.NewDecoder(&buf)
	if err := decoder.Decode(&decoded); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", decoded)
	// {ID:abc123 Amount:50}
}
```

`bytes.Buffer` is a standard library type that acts as an in-memory, growable byte buffer — here it stands in for "wherever the encoded bytes end up," whether that's eventually a network connection or a file. `gob.NewEncoder(&buf)` creates an encoder that writes gob-encoded bytes into `buf`, and `encoder.Encode(tx)` performs the actual serialization. `gob.NewDecoder(&buf)` and `decoder.Decode(&decoded)` reverse the process. Notice there's no struct tag needed at all here — gob works directly off of Go's own exported field names and types, with no separate annotation layer required, since it only ever needs to make sense to another Go program.

Gob has one behavior worth calling out explicitly: unlike JSON, a gob encoder and decoder are somewhat stateful across multiple `Encode`/`Decode` calls on the same stream — the first time a given type is encoded, gob transmits a compact description of that type's shape once, and subsequent values of the same type reuse that description rather than re-describing it every time. This makes gob particularly efficient when you're encoding *many* values of the same type in sequence (exactly GoChain's situation — many transactions, many blocks), at the cost of being a format that only Go itself can realistically read back.

---

## 4. Option Three: A Hand-Rolled Binary Format

The third option skips a general-purpose library format entirely and writes out exactly the bytes you want, in exactly the layout you choose, using `encoding/binary` — a low-level standard library package for converting between Go's built-in numeric types and raw bytes, in a specific byte order.

```go
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type toyTransaction struct {
	ID     string // we'll assume, for this sketch, IDs are exactly 8 bytes
	Amount int64
}

// encode writes a hand-rolled, fixed-layout binary representation:
// 8 bytes for the ID, followed by 8 bytes (big-endian) for the amount.
// This is the same general idea real blockchains use for their most
// performance- and size-sensitive on-disk and on-wire formats.
func encode(tx toyTransaction) []byte {
	buf := new(bytes.Buffer)
	buf.WriteString(tx.ID) // exactly 8 bytes, by our stated assumption
	binary.Write(buf, binary.BigEndian, tx.Amount)
	return buf.Bytes()
}

func decode(data []byte) toyTransaction {
	id := string(data[0:8])
	amount := int64(binary.BigEndian.Uint64(data[8:16]))
	return toyTransaction{ID: id, Amount: amount}
}

func main() {
	tx := toyTransaction{ID: "abcd1234", Amount: 50}

	encoded := encode(tx)
	fmt.Printf("encoded: %d bytes -> %x\n", len(encoded), encoded)
	// encoded: 16 bytes -> 6162636431323334000000000000032

	decoded := decode(encoded)
	fmt.Printf("%+v\n", decoded)
	// {ID:abcd1234 Amount:50}
}
```

`binary.Write(buf, binary.BigEndian, tx.Amount)` writes `tx.Amount` (an `int64`) as exactly 8 raw bytes into `buf`, in **big-endian** byte order — meaning the most significant byte comes first, matching the byte order convention most network protocols and many blockchain formats (including Bitcoin's block header, in most of its fields) use, specifically so that two different systems agree unambiguously on how to interpret the same bytes. `binary.BigEndian.Uint64(data[8:16])` does the reverse: reading exactly 8 bytes back out as a `uint64`, which we then convert to `int64`.

This approach produces the smallest possible output — exactly 16 bytes for this toy transaction, with zero overhead spent describing field names or types, since the *code itself* (both the encoder and the decoder) already encodes an agreed-upon, fixed layout. The cost is that you, the programmer, are now fully responsible for getting that layout exactly right on both ends — get the byte offsets wrong, or the byte order wrong, or forget to account for a variable-length field (like a transaction ID that isn't always exactly 8 bytes, which real transaction IDs, being full 32-byte hashes, actually are — this sketch simplified for the example), and you get silently wrong data rather than a clear encoding error.

---

## 5. Comparing the Three, and Why GoChain Starts With gob

```
                     encoding/json         encoding/gob          encoding/binary
                    -----------------    -----------------    ---------------------
Human-readable?          Yes                   No                     No
Cross-language?          Yes                   No (Go only)           Yes (if you
                                                                        document the
                                                                        layout)
Output size              Larger                Smaller                Smallest
                     (field names          (compact,              (no metadata
                      repeated as text)     type sent once)         at all)
Effort to use            Very low              Very low              Highest —
                                                                       you write every
                                                                       byte yourself
Used by real                                                        Bitcoin, Ethereum,
blockchains for                                                     and most production
core data?               Rarely                Never                blockchains — yes
```

Given this comparison, why does GoChain start with `encoding/gob` for its early chapters, rather than jumping straight to a hand-rolled binary format the way a "real" production blockchain would?

The honest answer is a deliberate trade-off between **learning pace** and **production fidelity**, made explicitly rather than left implicit. In these early volumes, you are still learning what a block and a transaction even *are*, how hashing and signing work, and how a chain fits together — spending significant time getting exact byte offsets and byte-order details correct in a hand-rolled format would be effort spent on a genuinely separate skill, at a moment when it would compete for attention with the core concepts this course is building toward. `encoding/gob` gets you a working, correct, reasonably efficient serialization with almost no code, letting Volumes 2 and 3 focus entirely on hashing, signing, and chaining blocks together.

This is explicitly a decision GoChain revisits later, not a permanent one: once performance and cross-language compatibility genuinely start to matter — particularly once Volume 7's networking layer needs a wire format that could, in principle, be understood by a non-Go implementation of the same protocol, and once Volume 8's storage layer starts caring about exact on-disk size at scale — later chapters introduce a hand-rolled binary format for those specific, performance-sensitive paths, following the same `encoding/binary` approach sketched in Section 4. You are not learning a "wrong" approach now that gets thrown away; you're learning the right approach *for this stage of the project*, with an explicit, honest plan for when and why it changes.

---

## 6. Go's Testing Framework: go test and testing.T

From this chapter forward, every GoChain package ships with tests, written using Go's built-in testing framework — no external testing library required. A Go test file is an ordinary `.go` file, in the same package as the code it tests, whose filename ends in `_test.go`, containing functions with a specific signature.

```go
// hash_test.go — lives alongside hash.go in the same package.
package toy

import "testing"

func TestAddOne(t *testing.T) {
	result := AddOne(5)
	if result != 6 {
		t.Errorf("AddOne(5) = %d; want 6", result)
	}
}
```

A test function must be named starting with `Test`, take a single parameter `t *testing.T`, and return nothing. `*testing.T` is the type Go's testing framework passes in to let a test report failures — `t.Errorf(...)` records a failure with a formatted message *and* allows the rest of the test function to keep running afterward (useful if you want to check several things in one test and see all the failures at once, not just the first). `t.Fatalf(...)` is a related method that also records a failure but immediately stops that test function from continuing — useful when a later check would be meaningless or would panic if an earlier one failed (for instance, if a setup step that should return a valid pointer instead returns `nil`).

Running `go test ./...` (introduced already in Chapter 06's Makefile) compiles and runs every `_test.go` file across every package in the module, reporting `PASS` or `FAIL` for each one:

```bash
go test ./...
```

```
ok      github.com/you/gochain/core       0.004s
ok      github.com/you/gochain/crypto     0.002s
--- FAIL: TestAddOne (0.00s)
    hash_test.go:8: AddOne(5) = 7; want 6
FAIL
FAIL    github.com/you/gochain/toy        0.003s
```

Test files are never included in a normal `go build` or in the final compiled binary — they only run when explicitly invoked through `go test`, so there's no runtime cost or size penalty to shipping thorough tests alongside every package.

---

## 7. Table-Driven Tests

A **table-driven test** is a Go idiom for testing many input/output cases with one small block of test logic, rather than writing a separate, nearly-identical test function for every case. Instead of `TestAddOneWithFive`, `TestAddOneWithZero`, `TestAddOneWithNegative`, and so on, you define a slice (a "table") of cases and loop over it once.

```go
package toy

import "testing"

func TestAddOneTableDriven(t *testing.T) {
	// Each entry in this slice is one test case: an input and its
	// expected output, plus a short name describing what it exercises.
	cases := []struct {
		name  string
		input int
		want  int
	}{
		{name: "positive number", input: 5, want: 6},
		{name: "zero", input: 0, want: 1},
		{name: "negative number", input: -1, want: 0},
		{name: "large number", input: 999_999, want: 1_000_000},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := AddOne(tc.input)
			if got != tc.want {
				t.Errorf("AddOne(%d) = %d; want %d", tc.input, got, tc.want)
			}
		})
	}
}
```

`cases := []struct{ ... }{ ... }` declares an anonymous struct type inline (a struct type with no separate name, defined right where it's used) and immediately builds a slice of it, one entry per test case. `t.Run(tc.name, func(t *testing.T) { ... })` runs the given function as a named **subtest** — this means `go test -v` reports each case individually (`TestAddOneTableDriven/positive_number`, `TestAddOneTableDriven/zero`, and so on), and critically, if one case fails, the rest still run and report their own results too, rather than the whole test function stopping at the first failure.

This pattern will be used constantly starting in Volume 2 — testing `crypto.Hash` against a table of known input/expected-output pairs, testing `core.Block` validation against a table of valid and deliberately invalid blocks, and so on. Learning this idiom now, on a trivial example, means it's already second nature by the time it matters for real cryptographic correctness.

---

## 8. A Full Worked Test Suite

Let's bring everything in this chapter together: a small package with a real function, a serialization round-trip, and a complete table-driven test file exercising both.

```go
// toy/transaction.go
package toy

import (
	"bytes"
	"encoding/gob"
)

// Transaction is a deliberately simplified stand-in for core.Transaction,
// used here purely to demonstrate serialization and testing patterns
// before the real type exists starting in Volume 3.
type Transaction struct {
	ID     string
	Amount int64 // amount in gochips
}

// Serialize encodes a Transaction using gob, matching GoChain's early-
// volume choice of serialization format, explained in Section 5.
func (tx Transaction) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(tx); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DeserializeTransaction reverses Serialize, reconstructing a Transaction
// from its gob-encoded bytes.
func DeserializeTransaction(data []byte) (Transaction, error) {
	var tx Transaction
	buf := bytes.NewBuffer(data)
	if err := gob.NewDecoder(buf).Decode(&tx); err != nil {
		return Transaction{}, err
	}
	return tx, nil
}
```

```go
// toy/transaction_test.go
package toy

import "testing"

func TestTransactionSerializeRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		tx   Transaction
	}{
		{name: "typical transaction", tx: Transaction{ID: "tx-001", Amount: 50}},
		{name: "zero amount", tx: Transaction{ID: "tx-002", Amount: 0}},
		{name: "empty ID", tx: Transaction{ID: "", Amount: 100}},
		{name: "large amount", tx: Transaction{ID: "tx-003", Amount: 1_000_000_000}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := tc.tx.Serialize()
			if err != nil {
				t.Fatalf("Serialize() returned unexpected error: %v", err)
			}
			if len(encoded) == 0 {
				t.Fatalf("Serialize() returned empty bytes for %+v", tc.tx)
			}

			decoded, err := DeserializeTransaction(encoded)
			if err != nil {
				t.Fatalf("DeserializeTransaction() returned unexpected error: %v", err)
			}

			if decoded != tc.tx {
				t.Errorf("round trip mismatch: got %+v, want %+v", decoded, tc.tx)
			}
		})
	}
}

func TestDeserializeTransactionRejectsGarbage(t *testing.T) {
	garbage := []byte{0xFF, 0x00, 0x01, 0x02} // not valid gob data at all
	_, err := DeserializeTransaction(garbage)
	if err == nil {
		t.Fatal("expected an error decoding garbage bytes, got nil")
	}
}
```

Notice the structure: `TestTransactionSerializeRoundTrip` uses `t.Fatalf` for the setup steps (if `Serialize` itself fails, or returns something obviously wrong, there's no point continuing that particular subtest), but `t.Errorf` for the actual assertion being tested (the round-trip equality check) — this way a *failed comparison* still lets the test framework report it cleanly as a failed expectation, and other subtests in the table still run independently regardless of what happens in this one. `decoded != tc.tx` works here because `Transaction` contains only comparable fields (a `string` and an `int64`) — Go structs support the `==`/`!=` operators directly as long as every field is itself comparable, which won't always be true for GoChain's real types once they contain slices like `[]byte` (slices are not directly comparable with `==` in Go, a detail Volume 2 and 3's real tests will need to work around using `bytes.Equal` or `reflect.DeepEqual` instead).

`TestDeserializeTransactionRejectsGarbage` demonstrates an equally important kind of test: not just "does the happy path work," but "does invalid input fail the way we expect, rather than silently producing wrong data or crashing." Every GoChain package from here on is expected to include both kinds of tests — correctness on valid input, and graceful failure on invalid input — exactly the discipline Volume 2's cryptography code and Volume 3's block validation code will lean on heavily.

---

## Summary

- Structs must become raw bytes to be hashed, signed, sent over a network, or written to disk — this conversion is called serialization, and its reverse is deserialization.
- Serialization must be **canonical**: the same logical value must always produce identical bytes, or hashing (and therefore tamper-evidence) breaks.
- `encoding/json` is human-readable and cross-language but verbose; ideal for GoChain's future API layer (Volume 10), not for high-volume internal data.
- `encoding/gob` is Go-native, compact, and effortless to use, but only readable by other Go programs — GoChain's choice for its early chapters.
- A hand-rolled format built on `encoding/binary` gives the smallest size and full control over layout (matching how real blockchains like Bitcoin encode their data), at the cost of manual, error-prone byte-offset bookkeeping — GoChain adopts this later, once performance and cross-language wire compatibility genuinely matter.
- `go test` runs any function named `TestXxx(t *testing.T)` in a `_test.go` file; `t.Errorf` records a failure and continues, `t.Fatalf` records a failure and stops that test immediately.
- **Table-driven tests** define a slice of named input/expected-output cases and loop over them with `t.Run(name, func(t *testing.T) { ... })`, so many cases share one block of test logic and each reports its own pass/fail independently.
- A solid test suite checks both the happy path (valid data round-trips correctly) and the failure path (invalid data produces a clear error rather than silently wrong output or a crash) — a discipline every GoChain package follows from here on.

---

## Exercises

### Easy

1. **Take the `toyTransaction` JSON example from Section 2** and add a third field, `Timestamp int64`, with an appropriate `json` struct tag. Marshal a sample value and print the resulting JSON string, then unmarshal it back and confirm all three fields survive the round trip.

2. **Write a table-driven test (following Section 7's pattern) for a simple `Double(n int) int` function** that you write yourself (it should just return `n * 2`), covering at least four cases: a positive number, zero, a negative number, and a large number. Run `go test -v` and confirm each named subtest is reported individually.

3. **Using the gob example from Section 3**, measure (with `len(buf.Bytes())` or equivalent) the encoded size of the same `toyTransaction` value once using `encoding/json` and once using `encoding/gob`. Report both sizes in bytes and explain, in one or two sentences, which is smaller and why, referencing the comparison table in Section 5.

### Medium

4. **Extend the hand-rolled binary encoder/decoder from Section 4** to correctly support a variable-length `ID` string, rather than assuming exactly 8 bytes. Do this by first writing a 4-byte (`uint32`) length prefix for the ID, followed by the ID's actual bytes, followed by the 8-byte amount — and update `decode` to read the length prefix first before knowing how many bytes to read for the ID itself. Write a table-driven test confirming round-trip correctness for IDs of several different lengths, including an empty ID.

5. **Write `TestTransactionSerializeRoundTrip`-style tests (from Section 8) for a small `Block`-like toy struct** you define yourself, containing a `Height int64`, a `PrevHash []byte`, and a `Transactions []Transaction` (reusing the `Transaction` type from Section 8). Since `[]byte` and slices in general aren't directly comparable with `!=`, research and use `reflect.DeepEqual` (from the standard `reflect` package) instead for your round-trip comparison, and explain in a comment why `==`/`!=` wouldn't have compiled at all for this struct.

6. **Benchmark (using `time.Now()`/`time.Since()` around a loop, not yet Go's formal benchmarking framework) encoding 10,000 `Transaction` values (from Section 8) once with `encoding/json` and once with `encoding/gob`**, timing only the encoding step, not the setup. Report both timings and the total byte sizes produced by each, and write two or three sentences explaining whether the results match what Section 5's comparison table predicted.

### Hard

7. **Design and implement a canonical serialization test**: write a function that constructs the "same" logical `Transaction` value in two different ways that a careless implementation might accidentally serialize differently (for instance, if your struct had a `map[string]int64` field instead of two independent fields, since Go map iteration order is deliberately randomized) and write a test that would catch such non-determinism if it existed. Since the `Transaction` type from Section 8 doesn't actually have this problem, explain in a written paragraph exactly what change to the struct *would* introduce it, and why gob's handling of maps (research this specifically) does or does not protect you automatically.

8. **Implement Go's formal benchmarking framework properly**: write a `BenchmarkTransactionSerializeGob(b *testing.B)` function (following the `func BenchmarkXxx(b *testing.B)` signature and the `for i := 0; i < b.N; i++ { ... }` loop pattern — research this from Go's `testing` package documentation if you haven't seen it before) that benchmarks `Transaction.Serialize()` from Section 8, and run it with `go test -bench=. -benchmem`. Report the nanoseconds-per-operation and bytes-allocated-per-operation figures it prints, and explain what `-benchmem` adds over a plain `-bench=.` run.

9. **Build a small "format migration" tool**: write a function `MigrateJSONToGob(jsonData []byte) ([]byte, error)` that reads a `Transaction` encoded as JSON and re-encodes it as gob, plus the reverse function `MigrateGobToJSON`. Write a table-driven test suite that constructs several `Transaction` values, round-trips each one through both migration directions (JSON to gob to JSON, and gob to JSON to gob), and asserts the final value exactly matches the original at every step. Write a short paragraph explaining why a real GoChain migration — say, if the project ever needed to convert an old JSON-based data store to the gob-based format described in this chapter, or vice versa when eventually moving to the hand-rolled binary format from Section 4 — would need exactly this kind of round-trip testing to be trustworthy before running against real, valuable on-chain data.
