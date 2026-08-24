# Chapter 03: Environment Setup and the GoChain Project Scaffold

This chapter gets your machine ready and creates the actual project you'll be building for the rest of the course. By the end, you'll have Go installed, an editor configured, Git tracking your work, and a real `gochain` module with its first runnable program — the seed that every later volume grows from.

## Table of Contents

1. [What You Need Installed](#1-what-you-need-installed)
2. [Installing Go](#2-installing-go)
3. [Choosing and Setting Up an Editor](#3-choosing-and-setting-up-an-editor)
4. [Installing and Configuring Git](#4-installing-and-configuring-git)
5. [What Is a Go Module?](#5-what-is-a-go-module)
6. [Creating the gochain Module](#6-creating-the-gochain-module)
7. [The Full Package Layout, Volume by Volume](#7-the-full-package-layout-volume-by-volume)
8. [Writing the First main.go](#8-writing-the-first-maingo)
9. [Running and Building GoChain for the First Time](#9-running-and-building-gochain-for-the-first-time)
10. [Your First Commit](#10-your-first-commit)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. What You Need Installed

Three things, and only three, are required before you can start writing GoChain code:

1. **The Go toolchain** — the compiler, standard library, and command-line tools (`go build`, `go test`, `go run`, and so on).
2. **A code editor** — any editor works, but this chapter recommends a setup with good Go support.
3. **Git** — the version control system you'll use to track every change to GoChain across the whole course.

That's it. No database server, no Docker, no cloud account — those show up much later (Volume 8 for storage, Volume 13 for deployment). Everything through the first several volumes runs entirely on your own laptop with nothing but these three tools.

---

## 2. Installing Go

Go is distributed as a free installer for macOS, Windows, and Linux, straight from the official Go website. The exact download step looks different per platform, so rather than reproducing a set of screenshots that go stale, here is what to verify once you're done, on any platform, from a terminal:

```bash
go version
```

You should see output resembling:

```
go version go1.23.0 darwin/arm64
```

The exact version number will drift over time as new Go releases ship — anything reasonably recent (Go 1.21 or later) works fine for every example in this course. The important part of that output is that a version string prints at all, meaning Go is correctly installed and on your system's `PATH` — the list of directories your shell searches when you type a command name.

If `go version` fails with a "command not found" style error, the Go binary directory isn't on your `PATH` yet. On macOS and Linux this usually means adding a line like the following to your shell's startup file (`~/.zshrc`, `~/.bashrc`, or equivalent):

```bash
export PATH=$PATH:/usr/local/go/bin
```

After editing that file, either open a new terminal window or run `source ~/.zshrc` (adjusting for your actual shell config file) to pick up the change immediately.

It's also worth checking where Go will look for your workspace and installed tools:

```bash
go env GOPATH
```

Modern Go (using modules, which we set up in Section 5) doesn't require you to place your code inside `GOPATH` the way very old Go tutorials describe — you can create a GoChain project anywhere on your filesystem you like, which is exactly what we'll do.

---

## 3. Choosing and Setting Up an Editor

Any text editor can technically write Go code, but a few give you enormously more help than a plain text editor would:

- **Visual Studio Code** with the official Go extension (search "Go" by the Go Team at Google in the Extensions panel) is the most common choice — it gives you inline error checking, jump-to-definition, auto-formatting on save, and integrated test running.
- **GoLand** (JetBrains) is a paid, Go-specific IDE with very deep tooling, if you prefer a dedicated IDE experience.
- **Vim/Neovim** with `gopls` (Go's official language server) configured through a plugin like `vim-go` or a Neovim LSP client works well for developers already comfortable in a terminal-based editor.

Whichever you choose, make sure it's wired up to **`gofmt`** (Go's standard code formatter, which ships with the Go toolchain itself) so your files are automatically reformatted to Go's standard style every time you save. This matters more in Go than in many languages: `gofmt` isn't just a style preference, it's close to a community-wide standard, and every code example in this course assumes `gofmt`-formatted code (tabs for indentation, specific brace placement, and so on).

You can also run formatting manually, at any time, from the terminal:

```bash
gofmt -l .        # list files that are not correctly formatted
gofmt -w .        # rewrite them in place to be correctly formatted
```

---

## 4. Installing and Configuring Git

Git tracks every change you make to GoChain over the course of this project — genuinely useful given how much the codebase grows between Volume 1 and Volume 13. If Git isn't already installed, install it through your platform's usual method (Xcode Command Line Tools on macOS, a package manager on Linux, the official installer on Windows), then confirm it:

```bash
git --version
```

Set your identity once, globally, so every commit is correctly attributed:

```bash
git config --global user.name "Your Name"
git config --global user.email "you@example.com"
```

We'll initialize a Git repository for the `gochain` project specifically in Section 6, right after creating the module.

---

## 5. What Is a Go Module?

A **module** is Go's unit of dependency versioning and distribution — think of it as "one project, plus a manifest describing exactly which external packages it depends on and at which versions." Every module has a unique **module path**, typically matching where its source code lives (for example, a path under `github.com/<username>/<repo>`), which other code uses to import it.

For this entire course, GoChain's module path is:

```
github.com/you/gochain
```

This path doesn't need to correspond to an actual, currently-existing GitHub repository for the code in this course to work — Go only requires the path to be *unique enough not to collide* with real packages you might import, and consistent throughout the project. If you do want to publish GoChain to your own GitHub account later, you would replace `you` with your actual GitHub username and update the module path accordingly (Chapter 06 discusses this in more depth).

A module's identity lives in a single file at its root: `go.mod`. This file is Go's manifest — it declares the module's path, the minimum Go version it requires, and every external dependency it needs, pinned to an exact version. You'll see this file grow slightly over the course as GoChain adds a handful of external dependencies (a CLI framework in Volume 10, an embedded database in Volume 8), but it starts nearly empty.

---

## 6. Creating the gochain Module

Time to create the real project. Pick a location on your filesystem for your course work — anywhere is fine — and run:

```bash
mkdir gochain
cd gochain
go mod init github.com/you/gochain
```

This produces a `go.mod` file that looks like:

```
module github.com/you/gochain

go 1.23
```

The `module` line declares the module path used in every import statement throughout the project — this is exactly why, from here on, any GoChain file that needs to import the (not-yet-written) crypto package will write `import "github.com/you/gochain/crypto"`, not a relative path. The `go` line declares the minimum Go language version this module requires; `go mod init` fills this in automatically based on the Go version installed on your machine when you ran the command.

Now initialize Git in this same directory:

```bash
git init
```

Before your first commit, it's worth adding a `.gitignore` file so compiled binaries and editor clutter never get accidentally committed:

```
# .gitignore
/gochain
/gochain-linux
*.exe
.DS_Store
data/
```

`data/` is included preemptively — starting in Volume 3, GoChain will write its blockchain data to a local `data/` directory, and that's exactly the kind of generated, machine-specific content that belongs in `.gitignore` rather than in version control.

---

## 7. The Full Package Layout, Volume by Volume

GoChain's code is organized into a small number of focused packages, each with one clear responsibility, mirroring the philosophy discussed in Chapter 02's tour of go-ethereum. Not every package is filled in yet — most start as an empty directory today and gain real code only once their volume "arrives." Here is the complete layout this course builds toward:

```
gochain/
├── go.mod
├── go.sum                  (appears once external deps are added)
├── main.go                 (temporary — moves into cmd/gochain later)
├── Makefile                (added in Chapter 06)
│
├── crypto/                 (Volume 2 — hashing, keys, signatures, Merkle trees)
│   ├── hash.go
│   ├── keys.go
│   └── merkle.go
│
├── core/                   (Volumes 3, 5 — the blockchain itself)
│   ├── block.go
│   ├── blockchain.go
│   ├── transaction.go
│   └── mempool.go
│
├── consensus/               (Volumes 4, 11 — proof of work, proof of stake)
│   ├── pow.go
│   └── pos.go
│
├── storage/                 (Volume 8 — persistence)
│   ├── store.go
│   └── bolt_store.go
│
├── wallet/                  (Volume 6 — keys, HD derivation, mnemonics)
│   ├── wallet.go
│   └── hd.go
│
├── network/                 (Volume 7 — P2P networking)
│   ├── node.go
│   ├── peer.go
│   └── protocol.go
│
├── vm/                      (Volume 9 — smart contracts)
│   ├── vm.go
│   ├── stack.go
│   └── opcodes.go
│
├── api/                     (Volume 10 — JSON-RPC, REST, WebSocket)
│   ├── rpc.go
│   ├── rest.go
│   └── websocket.go
│
└── cmd/
    └── gochain/             (Volume 10 onward — the CLI entrypoint)
        └── main.go
```

A few things worth noting about this layout up front, even though most of it won't matter until much later:

- **Each package name describes exactly one responsibility.** `core` never reaches into `network`'s internals, and `crypto` never needs to know anything about blocks or transactions at all — it just hashes bytes and signs bytes. This separation is what lets you swap, say, the storage engine in Volume 8 without touching consensus code.
- **`cmd/gochain` is deliberately separate from every library package.** It contains only the thin "glue" code that wires the library packages together into a runnable program — the actual logic always lives in `core`, `crypto`, `consensus`, and so on. Chapter 06 explains why this separation matters.
- **Dependencies flow in one direction.** `core` will import `crypto` (a block needs to hash itself). `network` will import `core` (a node needs to send and receive blocks). `crypto` will never import `core` — lower-level packages never depend on higher-level ones. Chapter 06 diagrams this dependency graph in full.

For today, you only need `main.go` at the project root — every directory above is a preview of what's coming, not something to create right now. Volume 1 chapters deliberately do not use any `core`, `crypto`, or other GoChain-specific types yet, since those packages don't exist until later volumes.

---

## 8. Writing the First main.go

Every GoChain milestone in this course starts from a program you can actually run and see output from — it's much more motivating than staring at code that doesn't do anything yet. Here is the very first version of `main.go`, placed directly at the root of the `gochain` module:

```go
// Command gochain is the entrypoint for the GoChain project — a real,
// working blockchain built from scratch across this course. This file
// starts out as a simple welcome banner and grows, chapter by chapter,
// into a full node, wallet, and command-line toolkit.
package main

import "fmt"

func main() {
	printBanner()
}

// printBanner prints a small welcome banner so the very first time you
// run this project, you get immediate, visible proof that your Go
// toolchain, module, and project layout are all wired up correctly.
func printBanner() {
	banner := `
   _____       _____ _           _        
  / ____|     / ____| |         (_)       
 | |  __  ___| |    | |__   __ _ _ _ __   
 | | |_ |/ _ \ |    | '_ \ / _' | | '_ \  
 | |__| | (_) | |____| | | | (_| | | | | |
  \_____|\___/ \_____|_| |_|\__,_|_|_| |_|

  GoChain — a real blockchain, built in Go, from scratch.
  Volume 1: Orientation & Go Essentials for Blockchain
`
	fmt.Println(banner)
}
```

A few small but deliberate choices here, worth calling out explicitly since they set patterns you'll see repeated throughout the course:

- `package main` declares this file as belonging to Go's special `main` package — the one package type that produces an executable program rather than an importable library. Every `cmd/...` directory in the layout from Section 7 will also be `package main` for exactly this reason.
- The `func main()` function is Go's designated entrypoint — when you run the compiled binary, execution starts here, and nowhere else.
- `printBanner` is pulled out as its own named function, rather than inlining the `fmt.Println` call directly inside `main`, even though it's only called once. This is a habit worth building early: as `main.go` grows across the course (it will eventually parse command-line flags, open a blockchain, start a network listener, and more) keeping `main` itself short and delegating each concern to a clearly named function keeps the entrypoint readable no matter how much the project grows.
- The comment directly above `package main` is a **package doc comment** — Go convention (enforced by tools like `go doc` and `golint`) is that the package's own documentation starts with the word "Command" (for executables) or the package's name (for libraries), immediately above the `package` declaration.

---

## 9. Running and Building GoChain for the First Time

With `go.mod` and `main.go` both in place, two commands matter right now, and you'll use both constantly for the rest of the course:

```bash
go run main.go
```

`go run` compiles and immediately executes a program in one step, without leaving a binary file behind afterward — perfect for quick iteration while developing. You should see the banner print to your terminal.

```bash
go build -o gochain .
./gochain
```

`go build` compiles the program into a standalone binary (named `gochain` here, via the `-o` flag) that you can run directly and distribute, without needing `go` installed on the machine that eventually runs it — exactly the deployment advantage discussed in Chapter 02, Section 5. The `.` tells `go build` to compile the package in the current directory.

```
   _____       _____ _           _        
  / ____|     / ____| |         (_)       
 | |  __  ___| |    | |__   __ _ _ _ __   
 | | |_ |/ _ \ |    | '_ \ / _' | | '_ \  
 | |__| | (_) | |____| | | | (_| | | | | |
  \_____|\___/ \_____|_| |_|\__,_|_|_| |_|

  GoChain — a real blockchain, built in Go, from scratch.
  Volume 1: Orientation & Go Essentials for Blockchain
```

If you see this output, your Go toolchain, module, and first program are all correctly wired together, and you're ready for every chapter that follows.

---

## 10. Your First Commit

With everything working, it's time to save this milestone in Git. Stage the two files that matter — `go.mod` and `main.go` — plus the `.gitignore` from Section 6:

```bash
git add go.mod main.go .gitignore
git commit -m "Initial GoChain scaffold: go.mod and welcome banner"
```

This is a genuinely meaningful commit, even though the code is simple: it's the seed every later chapter's code grows from. From here forward, this course expects you to commit at natural milestones — usually the end of each chapter, or after finishing a mini/major project — so you build the habit of a clean, reviewable history exactly the way a real production Go project would expect.

```
gochain/  (Git repository root)
├── .git/                   <- created by `git init`
├── .gitignore
├── go.mod
└── main.go

  git log --oneline
  --------------------------------
  a1b2c3d  Initial GoChain scaffold: go.mod and welcome banner
```

---

## Summary

- Only three tools are required to start: the Go toolchain, a code editor (ideally with Go language server support), and Git.
- `go version` confirms Go is installed and on your `PATH`; `gofmt` is Go's standard code formatter and should be wired into your editor to run on save.
- A **Go module** is a versioned unit of code identified by a **module path** (`github.com/you/gochain` for this course), declared in a `go.mod` file created by `go mod init`.
- GoChain's full package layout — `core`, `crypto`, `consensus`, `storage`, `wallet`, `network`, `vm`, `api`, and `cmd/gochain` — is laid out up front, even though most packages stay empty until their volume arrives later in the course.
- Each package has exactly one responsibility, dependencies flow in one direction (lower-level packages like `crypto` never import higher-level ones like `core`), and `cmd/` stays separate from library code — a pattern explained fully in Chapter 06.
- The first `main.go` prints a welcome banner, establishing the `package main` / `func main()` entrypoint pattern and the habit of pulling logic into small, named functions from the very first line of code.
- `go run` compiles and runs in one step for quick iteration; `go build -o gochain .` produces a standalone, distributable binary.
- The first Git commit — `go.mod`, `main.go`, and `.gitignore` — is the seed every later chapter's code builds on top of.

---

## Exercises

### Easy

1. **Run `go version` and `go env GOPATH GOMODCACHE`** on your machine, and write down all three outputs. Explain, in one or two sentences based on the documentation or `go help env`, what `GOMODCACHE` is used for.

2. **Modify `printBanner`** to also print the current date and a placeholder version string (for example, `"GoChain v0.1.0-dev"`) below the ASCII banner, using `time.Now()` from the standard `time` package. Run `go run main.go` again and confirm the new line appears correctly formatted.

3. **Create the full directory tree from Section 7** (every package directory, but leave each one empty except for a single placeholder file if your editor or Git won't track empty folders — a common trick is a `.gitkeep` file). Run `go build ./...` (the `./...` pattern tells Go to build every package in the module) and confirm it succeeds with no source files yet in the new directories.

### Medium

4. **Add a `-version` command-line flag** to `main.go` using the standard `flag` package, such that running `go run main.go -version` prints just a version string and exits immediately, without printing the full banner. Running `go run main.go` with no flags should still print the banner as before. (This previews the flag-parsing pattern used again in Chapter 21's chain inspector CLI.)

5. **Write a short (150-200 word) explanation** of what would go wrong, mechanically, if GoChain's `crypto` package tried to import `core` (the reverse of the intended dependency direction described in Section 7) — specifically, describe the compile-time error category this would eventually risk (an import cycle) even if it doesn't happen immediately with just these two packages, and explain why Go refuses to compile a true import cycle at all.

6. **Cross-compile a binary for a platform different from the one you're using**, following the pattern from Chapter 02, Section 5 (`GOOS=... GOARCH=... go build -o gochain-<platform> .`). Confirm the resulting binary's file type using your platform's `file` command (on macOS/Linux) or by checking its properties (on Windows), and write down what it reports — you should see it correctly identified as a binary for the *target* platform, not the one you built it on.

### Hard

7. **Research what `go.sum` is for** (you won't have one yet, since GoChain has no external dependencies as of this chapter) by reading `go help modules` or the official Go modules documentation, and write a 150-250 word explanation of the difference between what `go.mod` records and what `go.sum` records, and why a team working on the same project needs both files committed to version control, not just `go.mod`.

8. **Design (on paper, not in code) an alternate package layout** for GoChain that merges `consensus` directly into `core` instead of keeping it separate, and write a 200-300 word argument for why the course's actual design (keeping them separate) is likely to pay off later — specifically referencing Volume 11's plan to add `consensus.ProofOfStake` as a second, swappable consensus engine alongside `consensus.ProofOfWork`, and what would have to change in a merged design to support that.

9. **Set up a minimal CI check** (using GitHub Actions, if you have a GitHub account, or describe the exact steps in detail if you don't want to set one up yet) that runs `go build ./...` and `go vet ./...` automatically on every push to your GoChain repository. Write out the full YAML workflow file (or the equivalent step-by-step description) and explain, in your own words, what class of mistakes `go vet` catches that a successful `go build` alone would not.
