# Chapter 06: Project Layout and Go Modules

GoChain is going to grow into a dozen interconnected packages over the course of this book, and getting the shape of that growth right from the start avoids painful rewrites later. This chapter explains how Go modules version and share code, how GoChain's own packages will depend on each other without tangling into circular knots, and sets up the `Makefile` you'll run for the rest of the course.

## Table of Contents

1. [Recap: What a Module Is](#1-recap-what-a-module-is)
2. [How Packages Import Each Other](#2-how-packages-import-each-other)
3. [GoChain's Dependency Graph](#3-gochains-dependency-graph)
4. [Why Import Cycles Are Forbidden, and How to Avoid Them](#4-why-import-cycles-are-forbidden-and-how-to-avoid-them)
5. [Internal Packages](#5-internal-packages)
6. [Why cmd/ Stays Separate From Library Code](#6-why-cmd-stays-separate-from-library-code)
7. [Adding an External Dependency](#7-adding-an-external-dependency)
8. [Writing GoChain's Makefile](#8-writing-gochains-makefile)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Recap: What a Module Is

Chapter 03 introduced the basics: a **module** is Go's unit of dependency versioning and distribution, identified by a **module path** (`github.com/you/gochain` for this course) and declared in a `go.mod` file. Everything inside the `gochain/` directory tree — every package in the layout from Chapter 03, Section 7 — belongs to this one module.

A module can contain many **packages** — a package is simply a directory of `.go` files that all declare the same `package` name at the top and are compiled together as one unit. `gochain/core`, `gochain/crypto`, and `gochain/network` are each their own package, all living inside the single `gochain` module. This is an important distinction to keep straight: one module, many packages. Most real Go projects, including go-ethereum and GoChain, work this way — a single `go.mod` at the repository root, with dozens of packages organized into subdirectories underneath it.

---

## 2. How Packages Import Each Other

A Go package imports another package by its full **import path**, which is the module path plus the subdirectory the package lives in. For GoChain, that means:

```go
package core

import (
	"github.com/you/gochain/crypto" // imports the crypto package
)

func computeSomething(data []byte) []byte {
	return crypto.Hash(data) // calling an exported function from crypto
}
```

`import "github.com/you/gochain/crypto"` tells the Go compiler exactly where to find the `crypto` package: inside this same module, in the `crypto/` subdirectory. Once imported, any **exported** identifier from that package — in Go, exported means the name starts with an uppercase letter, like `Hash`, `KeyPair`, or `Sign` — becomes accessible as `crypto.Hash`, `crypto.KeyPair`, and so on. A lowercase name like `crypto.hashInternal` would *not* be accessible from outside the `crypto` package at all — Go uses capitalization itself as the visibility rule, with no separate `public`/`private` keyword needed.

This matters for API design: every GoChain package should expose a small, deliberate set of exported names (its actual public interface) while keeping implementation details unexported, so other packages can only depend on the parts that are meant to be depended on.

---

## 3. GoChain's Dependency Graph

Chapter 03 sketched GoChain's package layout. Now let's be explicit about which packages import which others, since this shape is what keeps the project maintainable as it grows across thirteen volumes. Here is the intended dependency graph for the whole project, with arrows meaning "imports":

```
                        cmd/gochain
                       /   |   |   \
                      /    |   |    \
                    api  network wallet
                      \    |   |    /
                       \   |   |   /
                        \  |   |  /
                         core (imports crypto, consensus)
                          |         \
                          |          \
                        crypto     consensus
                          |            |
                          +----(both import only the standard library)
                                       |
                                    storage
                          (imported directly by core AND network AND api)
                                       |
                                       vm
                          (imported by core, once contracts exist)
```

Reading this graph, a few rules become clear:

- **`crypto` sits at the bottom.** It has no dependency on any other GoChain package — it just hashes bytes and signs bytes using Go's standard library. Nothing about hashing or signing needs to know what a `Block` or a `Transaction` is.
- **`core` depends on `crypto`, `consensus`, and `storage`.** A `Block` needs to hash itself (`crypto`), a `Blockchain` needs to mine new blocks (`consensus`), and a `Blockchain` needs to persist itself (`storage`). But `core` never depends on `network`, `api`, or `wallet` — the blockchain data structure itself has no idea a network or an API even exists.
- **`network` depends on `core`.** A node needs to send and receive actual `Block` and `Transaction` values over the wire, so it must import `core` to know their shape. But `core` never imports `network` back — this asymmetry is exactly what keeps the graph a clean tree instead of a tangled web.
- **`wallet` depends on `crypto` (for keys and signing)** and, once transactions exist, on `core` (to build and submit `Transaction` values) — but neither `core` nor `crypto` ever imports `wallet` back.
- **`api`, `network`, and `wallet` are all consumed by `cmd/gochain`**, the actual executable, which sits at the very top of the graph and depends on nearly everything, directly or indirectly — but nothing depends on `cmd/gochain`, because it isn't a library at all (more on this in Section 6).

The single unifying rule across this whole graph: **dependencies only ever point downward, never upward, and never sideways-and-back**. `crypto` never needs to know `core` exists. `core` never needs to know `network` exists. Section 4 explains exactly why Go enforces this as a hard requirement, not just a style preference.

---

## 4. Why Import Cycles Are Forbidden, and How to Avoid Them

An **import cycle** happens when package A imports package B, and package B — directly, or through some chain of other packages — imports package A right back. Go's compiler refuses to build a program containing an import cycle at all; it's a hard compile error, not a warning:

```
import cycle not allowed
package github.com/you/gochain/core
	imports github.com/you/gochain/network
	imports github.com/you/gochain/core
```

Why does Go refuse this outright, rather than trying to figure out some sensible compilation order? Because a genuine cycle has no sensible order: to compile `core`, Go would need `network` to already be compiled (since `core` imports it), but to compile `network`, Go would need `core` to already be compiled (since `network` imports it) — there's no valid "first" package to start with. Some languages paper over this with lazy loading or forward declarations; Go's designers deliberately chose to disallow it entirely, on the theory that an import cycle is almost always a sign of a design problem — two concerns that should have been separated more cleanly, or split into a third, shared package both can depend on.

Suppose, hypothetically, `network` needed a helper function that felt like it belonged in `core`, and `core` also happened to need something `network` defined — a classic setup for an accidental cycle. The fix is almost always the same: extract the shared piece into its own package that both `core` and `network` can depend on without depending on each other.

```
   BEFORE (cycle — will not compile):        AFTER (extracted shared package):

   core -------> network                     core -------> shared
     ^              |                        network ----> shared
     |              |                        (core and network no longer
     +--------------+                         depend on each other at all)
```

GoChain's actual layout from Section 3 is deliberately designed to avoid ever needing this fix in the first place — `crypto` and `consensus` sit below `core`, and `core` sits below `network`, `wallet`, and `api`, with nothing ever pointing back down. As long as every new type and function you add respects this direction (ask yourself: "does this really belong in the lower-level package, or does it belong one level up?"), you'll never run into an import cycle across the whole course.

---

## 5. Internal Packages

Go has a special directory name with compiler-enforced meaning: any package living inside a directory literally named `internal/` can only be imported by code that lives inside the same parent tree as that `internal/` directory — external modules cannot import it at all, no matter how they try.

```
gochain/
├── core/
├── network/
│   ├── protocol.go
│   └── internal/
│       └── framing/          <- only importable from within gochain/network/...
│           └── framing.go
```

Code inside `gochain/network/internal/framing` can be imported by `gochain/network` itself, or by any other package nested underneath `gochain/network/`, but it *cannot* be imported by, say, `gochain/api`, or by a completely different project that happened to depend on the `gochain` module. This is useful for genuinely private implementation details you want to share across a few closely related files without accidentally making them part of your project's public API — for instance, GoChain's low-level wire-format framing logic (Chapter 44) might reasonably live in an `internal` package if it's only ever meant to be used by `network` itself, not by every other package in the project.

For most of this course, GoChain's packages are simple enough that `internal/` doesn't see heavy use — but it's worth knowing it exists and enforcing real, compiler-checked privacy, since production Go projects (including go-ethereum) use it extensively to draw a hard line between "our own implementation details" and "the API we're promising to keep stable for other people."

---

## 6. Why cmd/ Stays Separate From Library Code

`cmd/gochain/main.go`, referenced back in Chapter 03's layout, is intentionally the *only* place `package main` appears in the entire GoChain project (aside from the temporary root-level `main.go` from Chapter 03, which will eventually move here). Every other package — `core`, `crypto`, `consensus`, and so on — is an ordinary, importable library package.

This separation matters for a concrete reason: a `package main` file cannot be imported by anything else — it can only be compiled into a standalone executable. If GoChain's actual blockchain logic lived directly inside `cmd/gochain/main.go` instead of in `core`, nothing else could ever reuse it. Consider `gochain-wallet` and `gochain-explorer`, the two companion tools mentioned in this course's later volumes — both are entirely separate executables from `gochain` itself, but both need to work with the exact same `core.Transaction` and `core.Block` types, the exact same `crypto.Sign`/`crypto.Verify` functions, and so on. Because all of that logic lives in ordinary, importable packages rather than being buried inside one `main` function, `gochain-wallet` and `gochain-explorer` can simply `import "github.com/you/gochain/core"` and reuse it directly, with zero duplication.

```
gochain/
├── core/, crypto/, consensus/, ...    <- reusable library code, importable
│                                          by ANY executable in this module
└── cmd/
    ├── gochain/main.go                <- the main node binary
    ├── gochain-wallet/main.go          <- Volume 6: a separate wallet CLI,
    │                                      imports crypto, wallet, core
    └── gochain-explorer/main.go        <- Volume 10: a separate explorer
                                            binary, imports core, storage, api
```

The practical rule to hold onto for the rest of this course: **if you ever find yourself writing real logic — anything beyond parsing flags, wiring dependencies together, and calling into a library package — directly inside a `cmd/.../main.go` file, that's a signal the logic belongs in a proper package instead.** `main.go` should read almost like a short summary of "first do this, then do that," with the actual "how" living elsewhere.

---

## 7. Adding an External Dependency

So far, GoChain has depended on nothing but Go's own standard library. That changes starting in later volumes — Volume 8 adds an embedded database, Volume 10 adds the Cobra CLI framework. Adding an external dependency is a single command:

```bash
go get github.com/spf13/cobra@v1.8.0
```

`go get` downloads the specified module at the specified version, records it in `go.mod` under a `require` section, and — critically — also writes (or updates) a `go.sum` file, which stores cryptographic checksums of every dependency's exact contents. `go.sum` exists specifically so that anyone else building this project (including a CI server, or you on a different machine) can verify they're getting *exactly* the same code you tested against, not a tampered or accidentally different version — both `go.mod` and `go.sum` should always be committed to Git, never added to `.gitignore`.

```
go.mod                              go.sum
------------------------            ------------------------
module github.com/you/gochain       github.com/spf13/cobra v1.8.0 h1:...
                                     github.com/spf13/cobra v1.8.0/go.mod h1:...
go 1.23                             (one or more checksum lines per
                                     dependency, and per TRANSITIVE
require (                           dependency it pulls in)
	github.com/spf13/cobra v1.8.0
)
```

You won't need to run `go get` for anything in this volume — it's included here so the mechanism is familiar and unsurprising when Volume 8 and Volume 10 actually need it.

---

## 8. Writing GoChain's Makefile

A **Makefile** is a plain text file, read by the `make` command-line tool, that defines a set of named shortcuts (called **targets**) for common project tasks — instead of remembering and retyping a long `go build` invocation with a dozen flags, you type `make build`. `make` itself is a decades-old, extremely widely available build tool present on virtually every Unix-like system by default, and it works perfectly well as a lightweight task runner for a Go project, without needing anything Go-specific.

Here is the `Makefile` GoChain uses for the rest of this course, placed at the project root, right next to `go.mod`:

```makefile
# Makefile for the GoChain project.
# Run `make help` to see this list, or just run a target directly,
# e.g. `make build`, `make test`, `make run`, `make fmt`.

BINARY_NAME := gochain
CMD_DIR := ./cmd/gochain

.PHONY: build test run fmt vet clean help

## build: compile the gochain binary into ./bin/
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	go build -o bin/$(BINARY_NAME) $(CMD_DIR)

## test: run the full test suite, across every package, with verbose output
test:
	go test -v ./...

## run: build and immediately run the gochain binary
run: build
	./bin/$(BINARY_NAME)

## fmt: format every .go file in the project according to gofmt's rules
fmt:
	gofmt -w .

## vet: run go vet across every package to catch common mistakes
vet:
	go vet ./...

## clean: remove build artifacts
clean:
	rm -rf bin

## help: print this list of available targets
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
```

Walking through the parts worth explaining: `.PHONY: build test run fmt vet clean help` tells `make` that these target names are *actions* to perform, not files it should check the existence or timestamp of — without this, `make` would (following its original purpose as a tool for compiling C projects) assume a target named `build` refers to a file named `build`, and skip running it if such a file already happened to exist in the directory.

Each target is a name followed by a colon, with its shell commands indented on the following lines (traditionally with a literal tab character — most editors handle this automatically once a file is recognized as a Makefile). `build` depends on nothing and simply runs `go build`, placing the resulting binary in a `bin/` directory (which should be added to `.gitignore`, alongside the entries from Chapter 03). `run: build` declares that the `run` target depends on the `build` target — running `make run` first ensures `build` has been run (recompiling the binary if anything changed), and only then executes it. `test` runs `go test -v ./...`, where `./...` is Go's wildcard meaning "every package in this module, recursively" — Chapter 07 covers `go test` properly. `fmt` and `vet` wrap `gofmt` and `go vet` (a standard Go tool that catches common mistakes `gofmt` and the compiler themselves wouldn't, like a `Printf` call whose format string doesn't match its arguments).

With this file saved, the entire rest of the course reduces to a small, memorable set of commands:

```bash
make build   # compile the project
make test    # run every test
make run     # build and run
make fmt     # auto-format all code
```

Commit the `Makefile` now, alongside an update to `.gitignore` adding the `bin/` directory:

```bash
git add Makefile .gitignore
git commit -m "Add Makefile with build, test, run, and fmt targets"
```

---

## Summary

- One Go **module** (`github.com/you/gochain`, declared in `go.mod`) contains many **packages** — directories of `.go` files sharing a package name, imported by their full path (module path plus subdirectory).
- Capitalization is Go's visibility rule: exported (capitalized) names are accessible from other packages; unexported (lowercase) names are not.
- GoChain's dependency graph flows strictly downward: `crypto` and `consensus` depend on nothing project-specific; `core` depends on both; `network`, `wallet`, and `api` depend on `core`; `cmd/gochain` sits at the top and depends on everything it needs, directly or indirectly.
- **Import cycles** (A imports B imports A) are a hard compile error in Go, not just a warning — the fix is almost always extracting the shared piece both sides need into its own, lower-level package.
- **`internal/` packages** are compiler-enforced private implementation details, importable only from within their own parent directory tree, not from outside the module or from unrelated sibling packages.
- **`cmd/` stays separate from library code** so that real logic lives in reusable, importable packages (`core`, `crypto`, and so on), letting separate executables like `gochain-wallet` and `gochain-explorer` reuse GoChain's core logic without duplicating it.
- External dependencies are added with `go get module@version`, which updates both `go.mod` (the declared requirement) and `go.sum` (checksums for verification) — both files must always be committed to version control.
- GoChain's `Makefile` defines `build`, `test`, `run`, `fmt`, `vet`, and `clean` targets, giving the whole rest of the course a small, consistent set of commands to build, test, and run the project.

---

## Exercises

### Easy

1. **Draw (on paper or in a plain-text diagram, following the ASCII style used in Section 3) the dependency graph you'd expect once `wallet` is added**, given that `wallet` needs to generate key pairs (from `crypto`) and, eventually, build and sign real `core.Transaction` values. Confirm your diagram doesn't introduce any arrow pointing from `core` or `crypto` back toward `wallet`.

2. **Run `make help`** against the Makefile from Section 8 and confirm it lists all six targets with their one-line descriptions. Then add a new target, `make vet-and-test`, that runs `go vet ./...` followed by `go test ./...` in sequence (hint: you can depend on both existing targets, or just write both commands directly in the new target's body), and confirm it works.

3. **Deliberately introduce a two-package import cycle** in a small throwaway Go module (not inside GoChain itself — create a fresh, separate scratch module for this exercise): make package `foo` import package `bar`, and package `bar` import package `foo`. Run `go build ./...`, paste the exact error message Go produces, and identify which two lines in your two files you'd need to change to fix it.

### Medium

4. **Extend the `Makefile` from Section 8 with a `lint` target** that runs `gofmt -l .` (list-only mode, which prints filenames that are *not* correctly formatted without changing them) and fails the build (exits with a non-zero status) if any files are listed — research the shell trick needed to make a command's success depend on another command's output being empty (a common approach uses `test -z` combined with command substitution). Explain, in a comment above the target, why `gofmt -l` alone doesn't automatically fail a `make` target the way you want without this extra step.

5. **Design (on paper) where a hypothetical `gochain/testutil` package would sit in the dependency graph** if it existed to provide shared helper functions used only inside test files across multiple packages (for example, a function that builds a small, valid, fake `core.Blockchain` for tests to use). Which existing packages would `testutil` need to import, and is it safe for `core`'s own test files to import `testutil`, even though regular (non-test) `core` code never would? Justify your answer with reference to how Go treats `_test.go` files during a normal (non-test) build.

6. **Add a genuinely new dependency to a scratch Go module** (not GoChain itself) using `go get`, inspect the resulting `go.mod` and `go.sum` files, and write a short explanation (150-200 words) of what you observe: how many lines got added to `go.sum` versus `go.mod`, and why `go.sum` tends to have more entries than the single dependency you explicitly requested (hint: consider what a dependency's own dependencies are called, and how Go handles them).

### Hard

7. **Design a scenario (in writing, no code required) where GoChain's actual Volume 9 plans genuinely risk an import cycle**, based on the shared type contract at the start of this assignment: `vm` needs to execute opcodes that can check a transaction's signature (implying it needs something from `crypto` and possibly `core`), while `core`, once contracts exist, needs to be able to invoke the `vm` to execute a contract's code during block validation (implying `core` needs something from `vm`). Propose a concrete restructuring (which package should own what, or whether a new package is needed) that avoids the cycle, and justify it using the same "dependencies only point downward" principle from Section 3.

8. **Write a small Go program (in a scratch module) with three packages, `a`, `b`, and `c`,** where `a` imports `b`, `b` imports `c`, and `c` does NOT import `a` — a valid three-level chain, not a cycle — and confirm it builds successfully with `go build ./...`. Then change just `c` to import `a`, creating an indirect cycle through three packages rather than two, and confirm Go's error message correctly reports the full cycle path through all three packages, not just two.

9. **Propose and justify (in a 250-350 word write-up) a plan for splitting GoChain's `cmd/` directory into `cmd/gochain`, `cmd/gochain-wallet`, and `cmd/gochain-explorer` as three genuinely separate `go build` targets sharing the same `go.mod`**, including exactly what each binary's `main.go` should import from the shared library packages, how you'd extend the Makefile from Section 8 to build all three with a single `make build-all` target (or three separate targets), and what would go wrong (in terms of the dependency graph from Section 3) if `gochain-wallet`'s code accidentally ended up needing something from `gochain-explorer`'s code directly, rather than through a shared library package.
