# Chapter 04: A Go Crash Course for Blockchain Developers

This chapter covers every piece of core Go syntax and semantics you need before the first real blockchain code appears in Volume 2 — variables, structs, methods, interfaces, error handling, slices, maps, and pointers. Every single concept is demonstrated through one running example: a toy `Envelope` type that previews the shape of GoChain's real `Block` and `Transaction` types, so nothing here feels like a disconnected generic tutorial.

## Table of Contents

1. [Variables and Types](#1-variables-and-types)
2. [Structs: Giving Data a Shape](#2-structs-giving-data-a-shape)
3. [Methods: Attaching Behavior to Data](#3-methods-attaching-behavior-to-data)
4. [Pointers: Why a Block Holds a Hash, Not a Copy](#4-pointers-why-a-block-holds-a-hash-not-a-copy)
5. [Slices: Growable Lists](#5-slices-growable-lists)
6. [Maps: Fast Lookups by Key](#6-maps-fast-lookups-by-key)
7. [Error Handling the Go Way](#7-error-handling-the-go-way)
8. [Interfaces: Designing for Swappable Pieces](#8-interfaces-designing-for-swappable-pieces)
9. [Putting It All Together: A Tiny Envelope Chain](#9-putting-it-all-together-a-tiny-envelope-chain)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Variables and Types

A **variable** is a named slot in memory that holds a value of a particular **type** — a category describing what kind of value it can hold (a whole number, some text, true/false, and so on). Go is **statically typed**, meaning every variable's type is fixed and checked by the compiler before the program runs (Chapter 02, Section 4 covered why this matters for financial software specifically).

```go
package main

import "fmt"

func main() {
	var height int64 = 0          // explicit type: a 64-bit integer
	timestamp := 1_735_000_000    // ":=" infers the type from the value (an int here)
	sealed := false               // a boolean: true or false
	label := "genesis envelope"   // a string: text

	fmt.Println(height, timestamp, sealed, label)
}
```

`var height int64 = 0` declares a variable named `height` with an explicit type, `int64` — a 64-bit signed integer, which is exactly the type GoChain's `core.Block.Height` field uses, since block heights only ever grow and never need to be negative in practice, but `int64` gives generous headroom without overflowing for an extremely long time. `timestamp := 1_735_000_000` uses Go's short variable declaration syntax, `:=`, which lets the compiler infer the type from the value on the right — here, a plain `int`. The underscores in `1_735_000_000` are just visual separators Go allows in numeric literals for readability; they don't change the value at all.

Common Go types you'll see constantly in GoChain code: `int64` (fixed-size integers, used for amounts and heights), `string` (text, like an address), `bool` (true/false, like whether a block passed validation), `byte` (an alias for `uint8`, an 8-bit unsigned integer — hashes and signatures in GoChain are represented as `[]byte`, a slice of these), and custom struct types, which we turn to next.

---

## 2. Structs: Giving Data a Shape

A **struct** is a custom type that bundles several named fields together into one value — Go's way of describing "a thing with these properties." If a blockchain didn't need structs, every block's timestamp, transactions, and hash would have to be tracked as separate, disconnected variables, which would be unmanageable the moment you have more than one block.

Let's define a small `Envelope` type — our stand-in, for this chapter only, for the real `core.Block` type that arrives in Volume 3. Think of it as a physical envelope: it holds a message, records when it was sealed, and (once sealed) carries a wax-seal fingerprint of its own contents, echoing the analogy from Chapter 01.

```go
package main

import "time"

// Envelope is a toy stand-in for what core.Block will become in Volume 3.
// It bundles a message with metadata about when and how it was sealed.
type Envelope struct {
	Message   string    // the contents being sealed
	SealedAt  time.Time // when Seal() was called
	SealHash  []byte    // the "wax seal" fingerprint, empty until sealed
	IsSealed  bool      // whether Seal() has been called yet
}
```

`type Envelope struct { ... }` declares a new type named `Envelope` with four fields, each with its own name and type. `Message string` holds the actual content. `SealedAt time.Time` uses the standard library's `time.Time` type to record a specific moment. `SealHash []byte` is a **slice** of bytes — Section 5 covers slices properly, but for now, think of `[]byte` as "a growable list of raw bytes," exactly the shape GoChain uses for every hash, signature, and public key throughout the entire course. `IsSealed bool` is a simple flag.

You create a value of a struct type using a **struct literal**:

```go
func main() {
	env := Envelope{
		Message: "Alice lent Bob 10 gochips",
	}
	// Fields not mentioned in the literal (SealedAt, SealHash, IsSealed)
	// get their type's "zero value" automatically: an empty time, a nil
	// slice, and false, respectively.
}
```

Notice that `env`'s `IsSealed` field defaults to `false` and `SealHash` defaults to `nil` (an empty/unset slice) automatically — Go always gives every field a sensible **zero value** if you don't specify one, rather than leaving it as uninitialized garbage the way some lower-level languages would.

---

## 3. Methods: Attaching Behavior to Data

A **method** is a function attached to a specific type, allowing you to write `env.Seal()` instead of a free-floating `seal(env)`. This is how GoChain's real types will expose their behavior — `block.ComputeHash()`, `tx.Sign()`, `chain.AddBlock()` — and it's worth understanding the syntax now on our small `Envelope` example.

```go
import (
	"crypto/sha256"
	"time"
)

// Seal computes a fingerprint of the envelope's message and marks it
// sealed, exactly like pressing a wax seal into an envelope in Chapter 01.
// The (e *Envelope) part is the "receiver" — it's what makes this function
// a method ON Envelope, callable as env.Seal().
func (e *Envelope) Seal() {
	hash := sha256.Sum256([]byte(e.Message))
	e.SealHash = hash[:]     // convert the fixed-size [32]byte array to a []byte slice
	e.SealedAt = time.Now()
	e.IsSealed = true
}
```

The part in parentheses right after `func`, `(e *Envelope)`, is called the **receiver** — it declares that `Seal` is a method belonging to the `Envelope` type, and that inside the method body, `e` refers to the specific `Envelope` value the method was called on. Because the receiver is `*Envelope` (a pointer to an `Envelope`, not a plain `Envelope`), calling `env.Seal()` actually modifies the original `env` variable in place — Section 4 explains exactly why that pointer is necessary here.

```go
func main() {
	env := Envelope{Message: "Alice lent Bob 10 gochips"}
	fmt.Println("sealed before:", env.IsSealed) // false

	env.Seal()

	fmt.Println("sealed after:", env.IsSealed)  // true
	fmt.Printf("seal hash: %x\n", env.SealHash) // hex-printed fingerprint
}
```

This is exactly the pattern `core.Block` will use for its own `ComputeHash()` method starting in Chapter 17 — a method with a pointer receiver that fills in the struct's own hash field based on its other fields.

---

## 4. Pointers: Why a Block Holds a Hash, Not a Copy

A **pointer** is a value that holds the memory address of another value, rather than holding a copy of that value itself. Go writes a pointer type as `*T` (a pointer to a `T`), gets the address of a value with `&value`, and reads through a pointer back to the underlying value with `*pointer`.

Why does this matter for a blockchain specifically? Consider the alternative: what if, instead of storing a *hash* of the previous block, each block stored an actual, full *copy* of the entire previous block? Every block would then contain a copy of every block before it, recursively — block 100 would physically contain nested copies of blocks 1 through 99 inside it. This would be enormously wasteful of memory and disk space, and utterly impractical once a chain has thousands or millions of blocks.

Instead, a block stores a small, fixed-size **hash** — a fingerprint, not a copy — of the previous block's contents. This is precisely why `core.Block.PrevBlockHash` is typed `[]byte` (effectively a small handful of bytes, 32 for SHA-256) rather than `*Block` or `Block`: a hash is cheap to store and compare, yet still lets you verify the previous block's contents haven't changed, by recomputing its hash and checking it matches.

Let's see the general pointer-vs-copy distinction directly with `Envelope`, since it's the same underlying mechanism:

```go
func sealCopy(e Envelope) {
	// e here is a COPY of whatever Envelope was passed in — sealing it
	// only affects this local copy, and the caller's original is untouched.
	e.IsSealed = true
}

func sealPointer(e *Envelope) {
	// e here is a POINTER to the caller's original Envelope — sealing it
	// affects the exact same Envelope the caller has.
	e.IsSealed = true
}

func main() {
	env := Envelope{Message: "test"}

	sealCopy(env)
	fmt.Println(env.IsSealed) // false — sealCopy only sealed its own copy

	sealPointer(&env)
	fmt.Println(env.IsSealed) // true — sealPointer modified the real env
}
```

`sealCopy(env)` passes `env` **by value** — Go copies the entire struct into the function's local `e` parameter, so any change made inside `sealCopy` is invisible to the caller once the function returns. `sealPointer(&env)` passes `&env`, the *address* of `env`, so the function's `e` parameter points at the exact same memory the caller's `env` occupies — changes made through `e` are changes to the caller's own value.

This is exactly why `Seal()` in Section 3 was defined with receiver `(e *Envelope)` rather than `(e Envelope)` — a plain-value receiver would only ever seal a throwaway copy, and `env.IsSealed` would stubbornly stay `false` no matter how many times you called `env.Seal()`. Any method that needs to *modify* the receiver must use a pointer receiver; methods that only *read* from the receiver (and don't need to change it) can technically use either, but GoChain will consistently use pointer receivers for its core types to avoid ever accidentally copying a large struct like `core.Block` unnecessarily.

```
   Value (copy) semantics:              Pointer semantics:

   caller's env  --copy-->  e            caller's env  <--address--  e
   (independent, changes                 (same memory, changes
    to e don't affect env)                to e DO affect env)
```

---

## 5. Slices: Growable Lists

A **slice** is Go's everyday tool for representing an ordered, growable list of values of the same type — `[]byte` (a list of bytes, used for hashes and signatures throughout GoChain), `[]string`, and eventually `[]*Transaction` (a block's list of transactions, matching the canonical `core.Block.Transactions` field). A slice is conceptually a view onto an underlying array that can grow as needed.

```go
func main() {
	var messages []string // a nil slice — no elements yet, length 0

	messages = append(messages, "Alice lent Bob 10 gochips")
	messages = append(messages, "Bob lent Carol 5 gochips")
	messages = append(messages, "Carol lent Alice 3 gochips")

	fmt.Println(len(messages)) // 3 — the number of elements currently held

	for i, m := range messages {
		fmt.Printf("line %d: %s\n", i, m)
	}
}
```

`append(messages, ...)` returns a slice with the new element added — importantly, you always reassign the result back (`messages = append(messages, ...)`), because `append` may need to allocate a new, larger underlying array once the old one runs out of room, in which case the old slice variable would otherwise still point at the smaller, outdated array. `len(messages)` returns the current number of elements. The `for i, m := range messages` loop is Go's standard way to iterate a slice, giving you both the index (`i`) and a copy of each element (`m`) on every iteration.

A slice of envelopes, previewing exactly how `core.Blockchain` will hold its blocks starting in Chapter 18:

```go
func main() {
	var chain []*Envelope // a slice of POINTERS to Envelope

	for _, msg := range []string{"first", "second", "third"} {
		env := &Envelope{Message: msg}
		env.Seal()
		chain = append(chain, env)
	}

	for _, env := range chain {
		fmt.Printf("%s -> %x\n", env.Message, env.SealHash)
	}
}
```

Notice `[]*Envelope` here, not `[]Envelope` — a slice of *pointers* to envelopes, rather than a slice of envelopes themselves. This matters for the same reason discussed in Section 4: with `[]*Envelope`, appending an envelope to the chain doesn't copy it, and any later code holding the same pointer sees the exact same underlying data. `core.Blockchain`'s block list follows this same `[]*Block`-style pattern once it's introduced in Volume 3.

---

## 6. Maps: Fast Lookups by Key

A **map** associates keys with values and gives you near-instant lookup by key, unlike a slice, which requires scanning through elements one by one to find something specific. Maps are written `map[KeyType]ValueType`.

```go
func main() {
	sealedEnvelopes := make(map[string][]byte) // message text -> seal hash

	env := &Envelope{Message: "Alice lent Bob 10 gochips"}
	env.Seal()
	sealedEnvelopes[env.Message] = env.SealHash

	hash, found := sealedEnvelopes["Alice lent Bob 10 gochips"]
	if found {
		fmt.Printf("found seal: %x\n", hash)
	} else {
		fmt.Println("no such envelope")
	}
}
```

`make(map[string][]byte)` creates an empty, ready-to-use map from `string` keys to `[]byte` values. The two-value form of a map lookup, `hash, found := sealedEnvelopes[...]`, is an important Go idiom: `found` is `true` if the key actually exists in the map, and `false` if it doesn't (in which case `hash` is just the zero value for `[]byte`, which is `nil`) — this lets you distinguish "the key exists and its value happens to be empty" from "the key was never in the map at all," which matters enormously once GoChain starts using maps for things like the UTXO set (Volume 5) or the mempool's "which transaction IDs have I already seen" tracking (Chapter 05 of this volume).

---

## 7. Error Handling the Go Way

Go has no exceptions for ordinary error conditions. Instead, functions that can fail simply return an extra value of type `error` (a built-in interface — Section 8 explains interfaces properly) alongside their normal return value, and callers are expected to check it immediately.

```go
import (
	"errors"
	"fmt"
)

// Unseal returns the envelope's message, but only if it has actually
// been sealed — otherwise it returns an error explaining why it failed.
func (e *Envelope) Unseal() (string, error) {
	if !e.IsSealed {
		return "", errors.New("envelope has not been sealed yet")
	}
	return e.Message, nil
}

func main() {
	env := &Envelope{Message: "Alice lent Bob 10 gochips"}

	msg, err := env.Unseal()
	if err != nil {
		fmt.Println("error:", err)
		return // stop here — don't try to use msg, it's meaningless
	}
	fmt.Println("message:", msg)
}
```

`func (e *Envelope) Unseal() (string, error)` declares that `Unseal` returns two values: a `string` and an `error`. By strong Go convention, the error is always the *last* return value, and `nil` (Go's "no value"/"absence" marker) means "no error occurred." The pattern `if err != nil { ... return }` — sometimes jokingly called Go's most common line of code — is how virtually every function call that can fail gets checked, immediately, right where the call happens, rather than being caught somewhere far away in a try/catch block wrapping a large chunk of unrelated code.

This matters enormously for a blockchain. Consider `chain.AddBlock(newBlock)` in a later volume: if that call can fail (an invalid hash, a bad signature, a malformed transaction), Go's convention forces the calling code to explicitly decide, right at the call site, what happens on failure — silently ignoring an error requires deliberately writing `_ = chain.AddBlock(newBlock)` (assigning the error to `_`, Go's "discard this" placeholder), which is a visibly deliberate act, not something that happens by simply forgetting a `catch` block.

---

## 8. Interfaces: Designing for Swappable Pieces

An **interface** in Go describes a set of methods a type must have — it says nothing about *how* those methods are implemented, only that they exist with a matching signature. Any type that happens to implement all of an interface's methods automatically satisfies that interface, with no explicit "implements" declaration required (this is called **structural typing**, sometimes described as duck typing with compile-time checking).

This is precisely the mechanism that will let GoChain swap its storage engine (Volume 8: an in-memory version for tests, a real BoltDB-backed version for production) or its consensus algorithm (Volume 11: proof of work vs. proof of stake) without rewriting the code that *uses* those pieces. Let's see the idea on a small scale first, with `Envelope`:

```go
// Sealer describes anything that can seal a message and report back
// its fingerprint. Envelope will satisfy this interface, but so could
// a completely different type later, without either type needing to
// know about the other.
type Sealer interface {
	Seal() []byte
}

// WaxSeal implements Sealer using our toy SHA-256-based approach.
type WaxSeal struct {
	Message string
}

func (w WaxSeal) Seal() []byte {
	hash := sha256.Sum256([]byte(w.Message))
	return hash[:]
}

// processAndPrint works with ANYTHING that satisfies Sealer — it has no
// idea WaxSeal exists specifically, only that whatever it's given has
// a Seal() []byte method.
func processAndPrint(s Sealer) {
	fingerprint := s.Seal()
	fmt.Printf("fingerprint: %x\n", fingerprint)
}

func main() {
	w := WaxSeal{Message: "Alice lent Bob 10 gochips"}
	processAndPrint(w) // WaxSeal satisfies Sealer automatically
}
```

`type Sealer interface { Seal() []byte }` declares that anything satisfying `Sealer` must have a method `Seal() []byte`. `WaxSeal` never writes anything like `implements Sealer` — Go checks, purely by looking at `WaxSeal`'s method set, that it happens to have a matching `Seal() []byte` method, and that alone is enough for `processAndPrint(w)` to compile and work.

Why does this matter so much for GoChain specifically? Picture `gochain/storage.Store`, arriving in Volume 8, as an interface with methods like `PutBlock` and `GetBlock`. Every other package that needs to save or load a block — `core`, `network`, `api` — will be written against the `storage.Store` *interface*, never against a specific database implementation directly. That means a `BoltDBStore` (the real, disk-backed implementation) and a hypothetical `MemoryStore` (useful for fast, disk-free tests) can both satisfy `storage.Store`, and every other package works with either one completely unchanged. Volume 11 does exactly the same trick for `consensus.ProofOfWork` and `consensus.ProofOfStake` sharing one consensus interface. This is the single biggest reason interfaces matter in a project like this one: they let large parts of the system evolve independently of each other.

---

## 9. Putting It All Together: A Tiny Envelope Chain

Let's close the chapter by combining every concept above — structs, methods, pointers, slices, error handling, and a small interface — into one small, complete, runnable program: a miniature version of Chapter 01's wax-sealed envelope chain, now as real Go code.

```go
package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

// Envelope bundles a message with a seal (a fingerprint of the message
// PLUS the previous envelope's seal), previewing core.Block from
// Volume 3: a bundle of data plus a hash of the block before it.
type Envelope struct {
	Message      string
	PrevSealHash []byte
	SealHash     []byte
}

// Seal computes this envelope's own fingerprint from its message AND
// the previous envelope's seal, exactly like the pencil-and-paper hash
// chain from Chapter 01 — mutating the receiver in place, which is why
// the receiver here is a pointer, not a plain Envelope.
func (e *Envelope) Seal() {
	combined := append([]byte(e.Message), e.PrevSealHash...)
	hash := sha256.Sum256(combined)
	e.SealHash = hash[:]
}

// EnvelopeChain is a slice of pointers to sealed envelopes, mirroring
// how core.Blockchain will hold its blocks starting in Volume 3.
type EnvelopeChain struct {
	Envelopes []*Envelope
}

// AddMessage seals a new envelope, linking it to the chain's current
// last envelope, and appends it to the chain.
func (c *EnvelopeChain) AddMessage(message string) *Envelope {
	var prevHash []byte
	if len(c.Envelopes) > 0 {
		prevHash = c.Envelopes[len(c.Envelopes)-1].SealHash
	}
	env := &Envelope{Message: message, PrevSealHash: prevHash}
	env.Seal()
	c.Envelopes = append(c.Envelopes, env)
	return env
}

// Validate walks the whole chain and confirms every envelope's stored
// seal still matches what recomputing it right now would produce —
// exactly the isChainValid pseudocode from Chapter 01, Section 6.
func (c *EnvelopeChain) Validate() error {
	for i, env := range c.Envelopes {
		var expectedPrevHash []byte
		if i > 0 {
			expectedPrevHash = c.Envelopes[i-1].SealHash
		}
		combined := append([]byte(env.Message), expectedPrevHash...)
		recomputed := sha256.Sum256(combined)

		if fmt.Sprintf("%x", recomputed[:]) != fmt.Sprintf("%x", env.SealHash) {
			return errors.New(fmt.Sprintf("envelope %d has been tampered with", i))
		}
	}
	return nil
}

func main() {
	chain := &EnvelopeChain{}
	chain.AddMessage("Alice lent Bob 10 gochips")
	chain.AddMessage("Bob lent Carol 5 gochips")
	chain.AddMessage("Carol lent Alice 3 gochips")

	if err := chain.Validate(); err != nil {
		fmt.Println("chain is INVALID:", err)
	} else {
		fmt.Println("chain is valid")
	}

	// Now tamper with an old envelope, exactly like Chapter 01's exercise.
	chain.Envelopes[0].Message = "Alice lent Bob 1000 gochips"

	if err := chain.Validate(); err != nil {
		fmt.Println("chain is INVALID:", err) // this branch now runs
	} else {
		fmt.Println("chain is valid")
	}
}
```

Running this program prints "chain is valid" first, then, after the tampering line, "chain is INVALID: envelope 0 has been tampered with" — the exact behavior predicted by hand in Chapter 01. Every piece of syntax in this program is something covered in this chapter: `Envelope` and `EnvelopeChain` are structs; `Seal`, `AddMessage`, and `Validate` are methods, each using a pointer receiver because they read or write shared state; `[]*Envelope` is a slice of pointers; `error` is returned from `Validate` following Go's standard error-handling convention. Nothing here is new syntax — it's the same handful of ideas, recombined into something that actually works.

---

## Summary

- Go is statically typed: every variable has a fixed type checked at compile time, using either `var name Type = value` or the shorthand `name := value`.
- A **struct** bundles named fields into one custom type (`Envelope` here, `core.Block` and `core.Transaction` starting in Volume 3); uninitialized fields get sensible zero values automatically.
- A **method** is a function attached to a type via a receiver (`func (e *Envelope) Seal()`); pointer receivers (`*Envelope`) let a method modify the original value rather than a throwaway copy.
- A **pointer** (`*T`, `&value`, `*pointer`) holds an address rather than a copy — this is exactly why a block stores a small hash of the previous block rather than an expensive full copy of it.
- **Slices** (`[]T`) are growable, ordered lists, built with `append`; **maps** (`map[K]V`) give near-instant lookup by key and support the two-value `v, ok := m[key]` pattern to detect a missing key.
- Go handles errors by returning an `error` as the last return value and checking it immediately with `if err != nil`, rather than using exceptions — this forces every failure path to be explicitly handled at the call site.
- **Interfaces** describe a required set of methods, satisfied automatically (structurally) by any type that happens to implement them — the exact mechanism that will let GoChain swap storage engines and consensus algorithms later without rewriting the code that uses them.
- The chapter's closing `EnvelopeChain` program combines all of the above into a small, complete, tamper-detecting hash chain — a working preview of `core.Block` and `core.Blockchain`, arriving for real in Volume 3.

---

## Exercises

### Easy

1. **Add a `CreatedBy string` field** to the `Envelope` struct from Section 2, update the struct literal example to set it, and print it out. Then explain, in one or two sentences, what zero value this field would get if you constructed an `Envelope{}` without mentioning `CreatedBy` at all.

2. **Write a method `func (e *Envelope) Summary() string`** that returns a one-line human-readable description of an envelope, such as `"sealed envelope: 'Alice lent Bob 10 gochips' (a1b2c3...)"` using the first few hex characters of `SealHash`. Call it from `main` and print the result.

3. **Given a `map[string]int64` called `balances`** mapping an address string to an integer balance, write a small program that looks up a specific address using the two-value form (`bal, ok := balances[addr]`), and prints `"unknown address"` if `ok` is `false`, or the balance otherwise. Test it with both an address that exists in the map and one that doesn't.

### Medium

4. **Modify the `EnvelopeChain.Validate()` method from Section 9** so that instead of returning as soon as it finds the *first* tampered envelope, it collects *every* tampered envelope's index into a `[]int` slice and returns that slice along with an error only if the slice is non-empty. Test it against a chain where you've tampered with two different envelopes, and confirm both indices are reported.

5. **Write an interface `Validator` with one method, `Validate() error`**, and make both `*Envelope` (add a simple validation rule of your choosing, such as "Message must not be empty") and `*EnvelopeChain` satisfy it. Write a function `runValidation(v Validator)` that calls `.Validate()` on anything satisfying the interface and prints either "OK" or the error, then call it with both an `*Envelope` and an `*EnvelopeChain`.

6. **Explain, with a runnable example**, what happens if you pass an `Envelope` (not a `*Envelope`) into a function parameter of type `Envelope`, mutate a field inside that function, and then check the field's value in the caller afterward. Then rewrite the same example using `*Envelope` instead, and show the difference in output with a comment explaining exactly why it differs, tying your explanation back to Section 4.

### Hard

7. **Design and implement a small `EnvelopeStore` interface** with methods `Save(e *Envelope) error` and `Load(sealHashHex string) (*Envelope, error)`, then implement two different types that satisfy it: an `InMemoryEnvelopeStore` (backed by a `map[string]*Envelope`) and a `LoggingEnvelopeStore` that wraps another `EnvelopeStore` and simply prints a message before and after every `Save`/`Load` call, delegating the actual work to the wrapped store. Demonstrate that code written against the `EnvelopeStore` interface works identically with either implementation. (This "wrapping" pattern previews how GoChain might add logging or caching around `storage.Store` later without changing its interface.)

8. **Benchmark, with real numbers you measure yourself** (using Go's `time` package around a loop, no need for the formal `testing` benchmark framework yet — that's covered in Chapter 07), the difference between looking up a specific message's seal hash by (a) scanning a `[]*Envelope` linearly with a loop comparing `Message` fields, versus (b) looking it up directly in a `map[string]*Envelope` keyed by message. Run both approaches against a chain of at least 100,000 envelopes and report the actual timing difference you observed, explaining in your own words why the gap should grow as the chain gets longer.

9. **Extend the `EnvelopeChain` from Section 9 to support "forks"**: allow `AddMessage` to optionally target *any* existing envelope as its "previous" envelope (not just the current last one), producing a tree of envelopes rather than a single linear chain. Write a function `LongestPath(chain *EnvelopeChain) []*Envelope` that finds and returns the longest single unbroken path from the genesis envelope (the one with a `nil` `PrevSealHash`) to any leaf. Explain, in a short paragraph, how this tree-of-envelopes structure and its "longest path" rule previews the fork-handling problem and the longest-chain rule covered properly in Chapter 50.
