# Chapter 24: Packages, Modules, and Go Workspace

Go's module system is how dependencies are versioned and distributed. This chapter covers how Go code is organized into packages and modules, how the module graph works, and how multi-module workspaces (`go.work`) solve the local development problem.

## Table of Contents

1. [Packages](#1-packages)
2. [Modules and go.mod](#2-modules-and-gomod)
3. [Versioning and the Module Graph](#3-versioning-and-the-module-graph)
4. [Internal Packages](#4-internal-packages)
5. [Go Workspace](#5-go-workspace)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. Packages

A **package** is a directory of `.go` files that share the same `package` declaration. It is the unit of compilation and the unit of encapsulation.

```
myapp/
├── go.mod
├── main.go               package main
├── store/
│   ├── store.go          package store
│   └── store_test.go     package store_test  (external test package)
└── store/memory/
    └── memory.go         package memory
```

### Import paths

The import path is the module path + the directory path from the module root:

```go
// go.mod says: module github.com/alice/myapp

import "github.com/alice/myapp/store"         // store package
import "github.com/alice/myapp/store/memory"  // memory subpackage
```

### Exported vs unexported

An identifier is exported if it starts with an **uppercase letter**. Everything else is package-private.

```go
package store

type Store struct { ... }       // exported — usable outside this package
type entry struct { ... }       // unexported — only accessible in this package

func New() *Store { ... }       // exported
func (s *Store) validate() { }  // unexported method
```

### Package-level initialization

`init()` functions run before `main()`, after all variable initializations in the package. Each file can have multiple `init()` functions.

```go
package db

var defaultPool *sql.DB

func init() {
    // init runs once at startup — use sparingly
    // Prefer explicit initialization over hidden init() side effects
}
```

### The blank identifier for side effects

```go
import _ "github.com/lib/pq"  // import for side effects only (registers pq driver)
```

---

## 2. Modules and go.mod

A **module** is a collection of related packages versioned together. It is defined by a `go.mod` file at its root.

```
module github.com/alice/myapp

go 1.22

require (
    github.com/go-chi/chi/v5 v5.0.12
    github.com/jmoiron/sqlx v1.3.5
)

require (
    // Indirect dependencies (transitive)
    github.com/lib/pq v1.10.9 // indirect
)
```

### go.sum — the lockfile

`go.sum` records the expected cryptographic hash of each module. Never hand-edit it. It is checked into version control.

```
github.com/go-chi/chi/v5 v5.0.12 h1:...
github.com/go-chi/chi/v5 v5.0.12/go.mod h1:...
```

### Essential module commands

```bash
# Add a dependency (updates go.mod and go.sum)
go get github.com/go-chi/chi/v5@latest

# Add a specific version
go get github.com/go-chi/chi/v5@v5.0.12

# Remove unused dependencies from go.mod
go mod tidy

# Download all dependencies to module cache
go mod download

# Show the current module dependency graph
go mod graph

# Verify that go.sum matches actual downloads
go mod verify

# Show where a package comes from
go list -m -json github.com/go-chi/chi/v5

# Replace a module with a local fork (useful for patching)
go mod edit -replace github.com/original/pkg=../local-fork
```

### Module cache

Downloaded modules go to `$GOPATH/pkg/mod` (usually `~/go/pkg/mod`). The cache is content-addressed and shared across all projects on your machine.

```bash
# Clear the module cache
go clean -modcache
```

---

## 3. Versioning and the Module Graph

Go modules follow **semantic versioning**: `v<major>.<minor>.<patch>`.

### Minimum version selection (MVS)

Go uses MVS — no floating ranges, no `^`, no `~`. Each requirement specifies a minimum version. When two packages require the same module at different versions, Go picks the **highest minimum** required.

```
myapp requires:
    lib-a @ v1.2.0
    lib-b @ v1.3.0

lib-a requires:
    lib-c @ v1.1.0

lib-b requires:
    lib-c @ v1.4.0

→ Go selects lib-c v1.4.0 (the highest minimum)
```

This is **deterministic** — the same `go.mod` always produces the same build.

### Major version suffixes

A major version bump (v2+) breaks the import path contract, so it becomes part of the import path:

```go
// v1
import "github.com/go-chi/chi"

// v5 — different import path
import "github.com/go-chi/chi/v5"
```

Both can be imported in the same program if needed (they're distinct packages).

### Upgrading and downgrading

```bash
# Upgrade to latest patch/minor (within major)
go get github.com/some/pkg@latest

# Upgrade to a specific version
go get github.com/some/pkg@v1.5.0

# Downgrade
go get github.com/some/pkg@v1.3.2

# After upgrading, tidy up
go mod tidy
```

### `replace` directive

Override a module with a local path or fork:

```
replace github.com/original/lib => ../my-fork
replace github.com/original/lib v1.2.3 => github.com/my-fork/lib v1.2.3-fix
```

---

## 4. Internal Packages

An `internal` package can only be imported by code within its parent directory tree.

```
myapp/
├── internal/
│   ├── auth/
│   │   └── auth.go       # importable only from myapp/...
│   └── db/
│       └── db.go
├── handler/
│   └── handler.go        # can import myapp/internal/auth ✓
└── main.go               # can import myapp/internal/auth ✓

# External packages CANNOT import myapp/internal/auth ✗
```

This is enforced by the Go compiler, not just convention.

**Why `internal`?** It lets you share code between packages in your module without making it part of your public API. You can refactor internals without worrying about breaking external users.

```go
// internal/auth/token.go
package auth

// generateToken is unexported, but the whole package is internal
func generateToken(userID int64) string { ... }

// ValidateToken is exported, but only to packages in this module
func ValidateToken(token string) (int64, error) { ... }
```

---

## 5. Go Workspace

The **workspace** (`go.work`) solves a common pain point: developing two related modules simultaneously. Without it, you'd have to publish one module and `go get` it every time you make a change, or use `replace` directives in every `go.mod`.

### Scenario

You're developing an app and a shared library at the same time:

```
dev/
├── myapp/        # module: github.com/alice/myapp
│   ├── go.mod
│   └── main.go
└── mylib/        # module: github.com/alice/mylib
    ├── go.mod
    └── lib.go
```

### Create a workspace

```bash
cd dev/
go work init ./myapp ./mylib
```

This creates `go.work`:

```
go 1.22

use (
    ./myapp
    ./mylib
)
```

Now changes to `mylib` are immediately visible to `myapp` without publishing or `replace` directives. The workspace file tells `go` to use the local versions.

### Workspace commands

```bash
# Add a module to an existing workspace
go work use ./another-module

# Remove a module from workspace (edit go.work manually, or:)
go work edit -dropuse ./another-module

# Sync workspace (updates go.sum files for all modules)
go work sync

# Show the workspace graph
go list -m -json all
```

### go.work vs replace

| | `go work` | `replace` |
|--|-----------|-----------|
| Scope | Your local env only (never committed) | In go.mod (can be committed) |
| Purpose | Multi-module local dev | Fork patching |
| Affects others | No (git-ignored) | Yes (if go.mod is committed) |

**Always add `go.work` to `.gitignore`** — it's a local developer convenience, not part of the project definition.

```gitignore
go.work
go.work.sum
```

### Disabling workspace mode

```bash
GOWORK=off go build ./...  # Ignore the workspace file
```

---

## Summary

| Concept | Key fact |
|---------|----------|
| Package | Directory of `.go` files; unit of compilation and encapsulation |
| Module | Collection of packages with a single `go.mod`; versioned together |
| Import path | `module-path/subdir/subdir` |
| Exported | Uppercase first letter; accessible outside the package |
| MVS | Go always picks the highest minimum version — deterministic, no floating ranges |
| v2+ modules | Import path changes: `pkg/v2` — both v1 and v2 can coexist |
| `internal/` | Package only importable within its parent directory tree |
| `go.work` | Multi-module workspace for local development; never committed |

---

## Exercises

### Easy
1. Create a module `github.com/you/calc` with a `math` package containing `Add`, `Sub`, `Mul`, `Div` functions. Add a `math/advanced` subpackage with `Sqrt` and `Pow`. Write a `main.go` that imports and uses both. Run `go mod tidy` and verify the resulting `go.mod`.
2. Create an `internal/validator` package inside an existing module. Export a `Validate(email string) error` function. Verify that a test in the same module can import it, but that writing a standalone program with a different module path cannot.
3. Add `github.com/stretchr/testify` to a module using `go get`. Write one test using `assert.Equal`. Run `go mod tidy` and look at the diff in `go.mod` — which new direct/indirect dependencies were added?

### Medium
4. Create a workspace with two modules: `myapi` and `mycore`. `mycore` exports a `Config` struct and a `Load()` function. `myapi` imports and uses it. Modify `mycore` and verify the change is immediately visible in `myapi` without any `go get` or `replace` directive.
5. Explore the module graph: run `go mod graph` in a project with several dependencies. Parse the output into a map and find which dependency has the most dependents. Then run `go list -m -json all` and find which module has the most transitive dependencies.
6. Practice MVS: create a module that directly requires `lib-a @ v1.1.0` and `lib-b @ v1.3.0`, where `lib-a` requires `lib-c @ v1.0.0` and `lib-b` requires `lib-c @ v1.2.0`. Run `go list -m all` and verify that `lib-c @ v1.2.0` is selected. Now change your direct dependency on `lib-c` to `v1.5.0` — what happens?

### Hard
7. Build a **dependency audit tool**: given a `go.sum` file, fetch the checksums from the Go checksum database (`sum.golang.org`) and verify that every entry in `go.sum` matches. Flag any module whose hash differs — this is a supply-chain attack detection primitive.
8. Create a **multi-module monorepo**: three modules (`api`, `worker`, `shared`) inside a single repository. Set up a `go.work` file so all three can import from `shared` locally. Write a Makefile that builds each module independently and a CI script that verifies `GOWORK=off go build ./...` passes for each module independently (i.e., they don't accidentally depend on each other via the workspace).
