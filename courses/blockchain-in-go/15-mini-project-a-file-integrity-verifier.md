# Chapter 15: Mini Project — A File Integrity Verifier

Volume 2 has spent five chapters building cryptographic machinery — hashing (Chapter 08-09), Merkle trees (Chapter 10), key pairs and signatures (Chapters 11-13), and human-readable addresses (Chapter 14) — entirely in service of a blockchain that does not exist yet. Before Volume 3 finally builds `core.Block`, this chapter takes a detour that is not really a detour at all: it puts the very first, simplest tool from this volume, `gochain/crypto.Hash`, to work on a completely ordinary, completely useful problem that has nothing to do with cryptocurrency. You will build `fileguard`, a small command-line tool that fingerprints every file in a folder and tells you, later, exactly what changed. It is the same core idea behind `git status`, behind antivirus integrity scanners, and — not coincidentally — behind the tamper-evident chain of blocks Chapter 16 begins building next.

## Table of Contents

1. [What We're Building: fileguard](#1-what-were-building-fileguard)
2. [Why File Integrity Checking Is a Real Problem](#2-why-file-integrity-checking-is-a-real-problem)
3. [Design: The Manifest and Two Commands](#3-design-the-manifest-and-two-commands)
4. [Walking a Directory Tree in Go](#4-walking-a-directory-tree-in-go)
5. [Fingerprinting Files with gochain/crypto.Hash](#5-fingerprinting-files-with-gochaincryptohash)
6. [Saving and Loading the Manifest as JSON](#6-saving-and-loading-the-manifest-as-json)
7. [Comparing Two Manifests: Added, Changed, Deleted](#7-comparing-two-manifests-added-changed-deleted)
8. [Mini Project: fileguard](#8-mini-project-fileguard)
9. [Sample Terminal Session — Catching a Tampered File](#9-sample-terminal-session--catching-a-tampered-file)
10. [Extending fileguard Further](#10-extending-fileguard-further)
11. [Where This Fits in GoChain](#11-where-this-fits-in-gochain)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. What We're Building: fileguard

`fileguard` is a two-command CLI tool:

```
fileguard scan  -dir <folder> -manifest <file>   # fingerprint every file, save the result
fileguard check -dir <folder> -manifest <file>   # compare the folder against a saved manifest
```

`scan` walks a directory, computes a SHA-256 fingerprint of every file it finds using `gochain/crypto.Hash` (Chapter 09), and saves the results — a **manifest** — to a JSON file. `check` re-walks the same directory later, recomputes every fingerprint, loads the previously saved manifest, and reports exactly three things: which files are **new** since the last scan, which existing files **changed**, and which files were **deleted**. Nothing about this tool touches a blockchain, a private key, or a transaction — it is Chapter 09's `Hash` function, alone, doing genuinely useful work.

---

## 2. Why File Integrity Checking Is a Real Problem

"Did any of these files change since I last checked?" sounds like a trivial question, but answering it reliably is exactly the kind of problem hashing was built to solve, and real tools you likely already use solve it this same way. `git status` answers a version of this question every time you run it, by comparing file hashes against what Git last recorded. Antivirus and endpoint-security software often maintains a database of known-good file hashes, precisely so a single unexpected change to a system file — one a virus or an intruder might have modified — gets flagged immediately, without needing to understand *what* changed, only *that* something did.

The naive way to answer "did this file change?" — compare file sizes, or compare "last modified" timestamps — is unreliable in ways that matter. A file's modification timestamp can be reset by copying it, restoring it from a backup, or by an attacker deliberately covering their tracks, all without the *content* changing at all — a false negative, missing a real change nothing should have missed. Conversely, a file's modification time updates even if you resave a text file with byte-for-byte identical content, or when an editor "helpfully" touches every file in a directory it opens — a false positive, when nothing meaningful actually changed. A hash-based fingerprint sidesteps both problems at once: Chapter 08's determinism guarantee means identical content always produces an identical hash regardless of timestamps, and the avalanche effect (Chapter 08, Section 3) means *any* content change, however small, produces a completely different hash — exactly the property this tool needs, and exactly the property Chapter 08 spent an entire chapter establishing you could trust.

---

## 3. Design: The Manifest and Two Commands

A **manifest** is the saved record `scan` produces: a snapshot mapping every file's path (relative to the folder being scanned) to its hash at the moment of scanning. GoChain's manifest is a simple JSON document:

```json
{
  "root": "./project",
  "files": {
    "README.md": "5b91956b36f77846ddaa49115ac394bdd4376e99ed802eebe71aee9335db29f3",
    "docs/spec.txt": "7e390ac0c42f8a2b8ae7d2212a67077f710a2f6dbdda4a18b04b899f68a1ad3e",
    "notes.txt": "d00a461ea6a379c7cf2060f65bdcd1adec086b88c41067468ff5454c34a1da08"
  }
}
```

Storing relative paths, not absolute ones, is a deliberate choice: it means a manifest built for `/home/alice/project` still works correctly if the folder is later moved, renamed, or copied somewhere else entirely — the manifest describes the *contents* of the folder, not its specific location on one specific machine.

The tool's whole design falls out of this one data structure. `scan` builds a fresh manifest and writes it to disk. `check` builds a fresh manifest of the folder's *current* state, loads the previously saved manifest from disk, and diffs the two: any relative path present in the new manifest but missing from the old one is **added**; any path present in both but with a different hash is **changed**; any path present in the old manifest but missing from the new one is **deleted**.

```
   scan (first run)              check (later run)
   ────────────────              ─────────────────

   walk folder                   walk folder
        │                             │
        ▼                             ▼
   hash every file          hash every file (freshly, right now)
        │                             │
        ▼                             ▼
   write manifest.json       load the OLD manifest.json
                                       │
                                       ▼
                              compare old vs. new, file by file
                                       │
                                       ▼
                              report: added / changed / deleted
```

---

## 4. Walking a Directory Tree in Go

Go's standard library provides `filepath.WalkDir` (added in Go 1.16, and preferred over the older `filepath.Walk` for performance reasons not important here) specifically for recursively visiting every file and subdirectory under a root path:

```go
// cmd/fileguard/manifest.go
package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func buildManifest(root string) (*Manifest, error) {
	m := &Manifest{Root: root, Files: make(map[string]string)}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil // only fingerprint files, not directories themselves
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		// Section 5 fills in the actual hashing here.
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}
```

`WalkDir` takes a root path and a callback function, and calls that callback once for every file and directory it finds, recursively, including the root itself. The callback receives the full path, a `fs.DirEntry` describing what was found (Go's lightweight alternative to a full `os.FileInfo`, sufficient for checking `d.IsDir()`), and an error if something went wrong reading that particular entry (a permissions problem, for instance) — returning that error from the callback stops the walk and propagates it up to `WalkDir`'s own return value, exactly the error-propagation discipline Chapter 13, Section 11 established for `crypto.Sign` and `crypto.Verify`. `filepath.Rel(root, path)` computes exactly the relative path Section 3 explained the manifest needs to store.

---

## 5. Fingerprinting Files with gochain/crypto.Hash

Filling in the hashing step is now almost anticlimactic, which is exactly the point of building `Hash` carefully back in Chapter 09 — the hard work of getting canonical, trustworthy hashing right happened once, then, and every later use of it, including this one, is just a function call:

```go
import (
	"encoding/hex"

	"github.com/you/gochain/crypto"
)

// (continuing the WalkDir callback from Section 4)

data, err := os.ReadFile(path)
if err != nil {
	return err
}

hash := crypto.Hash(data)
m.Files[rel] = hex.EncodeToString(hash)
```

`os.ReadFile` reads an entire file into memory as a `[]byte` — the exact input shape `crypto.Hash` expects, with no conversion needed. The result is stored as a hex string (Chapter 09, Section 3), because a JSON manifest is meant to be human-readable and diffable with an ordinary text tool, and raw binary hash bytes would not survive JSON encoding cleanly (JSON strings are text, not arbitrary bytes) the way a hex string does. This is the same trade-off Chapter 09 drew between `Hash` (for machines) and `HashHex` (for anything printed, logged, or — as here — saved as readable text).

One limitation worth naming honestly, in the same spirit as this course's other honest simplifications (Chapter 12, Section 3; Chapter 14, Section 10): `os.ReadFile` loads an entire file into memory at once, which is fine for source code, documents, and configuration files, but would be a poor choice for fingerprinting a multi-gigabyte video file on a memory-constrained machine. A production-grade version would stream the file through `sha256.New()` and `io.Copy` in fixed-size chunks instead of reading it all at once — Exercise 7 asks you to build exactly that improvement.

---

## 6. Saving and Loading the Manifest as JSON

Chapter 07 already compared `encoding/json`, `encoding/gob`, and hand-rolled binary encodings for GoChain's internal use, and chose `gob` for the blockchain's own data. A manifest is different: it is meant to be opened, read, and diffed by a human with a text editor if needed, which is exactly `encoding/json`'s strength over `gob`'s compact-but-opaque binary format. Because `Manifest.Files` is a plain `map[string]string` — strings as both keys and values, never the map-of-structs case Chapter 09, Section 6 warned about — `encoding/json` has no canonical-serialization pitfalls to worry about here: JSON object keys are written as plain text strings either way, and this manifest is a saved snapshot for a human and a later run of this same program to compare, not something ever hashed or fed back into a cryptographic function itself.

```go
// cmd/fileguard/manifest.go (continued)

type Manifest struct {
	Root  string            `json:"root"`
	Files map[string]string `json:"files"` // relative path -> hex hash
}

func saveManifest(m *Manifest, path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func loadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
```

`json.MarshalIndent` (rather than plain `json.Marshal`) pretty-prints the output with two-space indentation, purely so a manifest file opened in a text editor is readable at a glance — exactly the kind of small touch that matters for a tool meant to be inspected by a human, not just consumed by another program.

---

## 7. Comparing Two Manifests: Added, Changed, Deleted

The comparison itself is two simple map lookups, run over the two manifests' file listings:

```go
// cmd/fileguard/diff.go
package main

import "sort"

type DiffResult struct {
	Added   []string
	Changed []string
	Deleted []string
}

func diffManifests(old, newM *Manifest) DiffResult {
	var result DiffResult

	// Anything in the NEW manifest not present in the OLD one is
	// either brand new, or an existing file whose hash changed.
	for path, newHash := range newM.Files {
		oldHash, existed := old.Files[path]
		if !existed {
			result.Added = append(result.Added, path)
		} else if oldHash != newHash {
			result.Changed = append(result.Changed, path)
		}
	}

	// Anything in the OLD manifest no longer present in the NEW one
	// was removed since the last scan.
	for path := range old.Files {
		if _, stillExists := newM.Files[path]; !stillExists {
			result.Deleted = append(result.Deleted, path)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Changed)
	sort.Strings(result.Deleted)

	return result
}
```

The two loops are deliberately separate and run in opposite directions — the first walks the *new* manifest looking for paths the *old* one does not explain (new or changed files), and the second walks the *old* manifest looking for paths the *new* one no longer has (deleted files) — because a single loop over just one of the two maps could never detect both directions of change on its own. Sorting each result slice before returning is a small but important touch: map iteration order is randomized in Go (Chapter 09, Section 6, the exact same fact that caused Note's `gob`-encoding bug), so without an explicit sort, `fileguard`'s output order would be different, unhelpfully, on every single run, even when comparing the exact same two manifests.

---

## 8. Mini Project: fileguard

Assembling Sections 4 through 7 into one runnable command-line tool, with `flag`-based subcommands in the spirit of Chapter 21's chain inspector CLI:

```go
// cmd/fileguard/main.go
package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/you/gochain/crypto"
)

type Manifest struct {
	Root  string            `json:"root"`
	Files map[string]string `json:"files"`
}

func buildManifest(root string) (*Manifest, error) {
	m := &Manifest{Root: root, Files: make(map[string]string)}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		m.Files[rel] = hex.EncodeToString(crypto.Hash(data))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

func saveManifest(m *Manifest, path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func loadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

type DiffResult struct {
	Added   []string
	Changed []string
	Deleted []string
}

func diffManifests(old, newM *Manifest) DiffResult {
	var result DiffResult

	for path, newHash := range newM.Files {
		oldHash, existed := old.Files[path]
		if !existed {
			result.Added = append(result.Added, path)
		} else if oldHash != newHash {
			result.Changed = append(result.Changed, path)
		}
	}
	for path := range old.Files {
		if _, stillExists := newM.Files[path]; !stillExists {
			result.Deleted = append(result.Deleted, path)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Changed)
	sort.Strings(result.Deleted)
	return result
}

func printDiff(diff DiffResult) {
	if len(diff.Added) == 0 && len(diff.Changed) == 0 && len(diff.Deleted) == 0 {
		fmt.Println("fileguard: no changes detected")
		return
	}
	for _, f := range diff.Added {
		fmt.Println("added:  ", f)
	}
	for _, f := range diff.Changed {
		fmt.Println("changed:", f)
	}
	for _, f := range diff.Deleted {
		fmt.Println("deleted:", f)
	}
}

func main() {
	scanCmd := flag.NewFlagSet("scan", flag.ExitOnError)
	scanDir := scanCmd.String("dir", ".", "directory to scan")
	scanManifest := scanCmd.String("manifest", "fileguard.json", "manifest file to write")

	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	checkDir := checkCmd.String("dir", ".", "directory to check")
	checkManifest := checkCmd.String("manifest", "fileguard.json", "manifest file to compare against")

	if len(os.Args) < 2 {
		fmt.Println("usage: fileguard <scan|check> [flags]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "scan":
		scanCmd.Parse(os.Args[2:])
		m, err := buildManifest(*scanDir)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		if err := saveManifest(m, *scanManifest); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Printf("fileguard: fingerprinted %d files, saved to %s\n", len(m.Files), *scanManifest)

	case "check":
		checkCmd.Parse(os.Args[2:])
		oldManifest, err := loadManifest(*checkManifest)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		newManifest, err := buildManifest(*checkDir)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		printDiff(diffManifests(oldManifest, newManifest))

	default:
		fmt.Println("usage: fileguard <scan|check> [flags]")
		os.Exit(1)
	}
}
```

Build it exactly the way Chapter 06's `Makefile` builds every other GoChain binary:

```
go build -o fileguard ./cmd/fileguard
```

---

## 9. Sample Terminal Session — Catching a Tampered File

Here is a real, complete session, run exactly as shown, against a small three-file project:

```
$ ls project
README.md   docs/   notes.txt

$ ./fileguard scan -dir ./project -manifest fileguard.json
fileguard: fingerprinted 3 files, saved to fileguard.json

$ cat fileguard.json
{
  "root": "./project",
  "files": {
    "README.md": "5b91956b36f77846ddaa49115ac394bdd4376e99ed802eebe71aee9335db29f3",
    "docs/spec.txt": "7e390ac0c42f8a2b8ae7d2212a67077f710a2f6dbdda4a18b04b899f68a1ad3e",
    "notes.txt": "d00a461ea6a379c7cf2060f65bdcd1adec086b88c41067468ff5454c34a1da08"
  }
}

$ ./fileguard check -dir ./project -manifest fileguard.json
fileguard: no changes detected
```

So far, nothing has changed, and `check` correctly reports exactly that. Now, simulate exactly the kind of tampering `fileguard` exists to catch: someone edits `notes.txt` (adding an unauthorized line), adds a brand-new file, and deletes `README.md` entirely:

```
$ echo 'Budget approved: $50,000.' >> project/notes.txt
$ echo 'Draft' > project/docs/new_file.txt
$ rm project/README.md

$ ./fileguard check -dir ./project -manifest fileguard.json
added:   docs/new_file.txt
changed: notes.txt
deleted: README.md
```

Every single change is caught, correctly categorized, and reported by relative path — exactly Section 3's design, running for real. Notice that `fileguard` did not need to know or guess *what* changed inside `notes.txt` (the avalanche effect made that unnecessary, exactly as Chapter 08, Section 3 promised) — it only needed to notice that the file's fingerprint no longer matched what was recorded, which is a far simpler and far more reliable check than trying to inspect file contents directly for "suspicious" changes.

---

## 10. Extending fileguard Further

A handful of natural next steps are left as exercises rather than built here, precisely so you get hands-on practice extending a tool you already understand completely, rather than only ever reading finished code:

- **Ignoring files and folders** — a real project has a `.git` directory, build artifacts, and other files nobody wants fingerprinted; a `.fileguardignore` file (Exercise 7) would filter `buildManifest`'s walk.
- **Detecting renamed files** — right now, renaming `notes.txt` to `notes-2024.txt` reports as one deletion plus one addition, even though the *content* never changed; a smarter diff could notice a deleted path's hash reappearing under a new path and report a rename instead (Exercise 8).
- **Streaming large files** — Section 5 already flagged that `os.ReadFile` loads a whole file into memory; Exercise 7 (in the code sense) asks you to swap in `sha256.New()` plus `io.Copy` for files above some size threshold.

---

## 11. Where This Fits in GoChain

`fileguard` never touches `core.Block`, `core.Blockchain`, or any type this course has not built yet — and that is exactly the point of this mini project. It demonstrates that Chapter 09's `Hash` function is not merely blockchain infrastructure; it is a genuinely general-purpose tool for answering "has this data changed?" anywhere that question matters. Chapter 16 picks up almost the identical idea and gives it a name this course will use for the rest of its life: a block's own hash is `fileguard`'s file-fingerprint idea, applied to a bundle of transactions instead of a folder of files, and Chapter 18's chain of `PrevBlockHash` references is `fileguard`'s "compare against what was recorded last time" idea, applied continuously, one block after the next, forever. If `fileguard`'s `check` command felt satisfying to watch catch a tampered file, Volume 3's tamper-evident blockchain is that same feeling, generalized into a permanent, cumulative structure.

---

## Summary

- `fileguard` is a two-command CLI (`scan` and `check`) built entirely on top of `gochain/crypto.Hash` (Chapter 09), with no blockchain-specific code involved at all.
- File-size and modification-timestamp checks are unreliable for detecting real content changes; a hash-based fingerprint is reliable because of the exact two properties Chapter 08 built around hashing: determinism and the avalanche effect.
- A **manifest** is a saved snapshot mapping each file's relative path to its hash at scan time, stored as human-readable JSON via `encoding/json` rather than `gob`, since a manifest is meant to be inspected directly.
- `filepath.WalkDir` recursively visits every file under a root directory; `os.ReadFile` plus `crypto.Hash` fingerprints each one; storing *relative* paths keeps a manifest valid even if the scanned folder is later moved.
- Comparing two manifests for added, changed, and deleted files requires two separate passes — one over the new manifest looking for unexplained paths, one over the old manifest looking for paths that disappeared — and each result should be sorted, since Go's map iteration order is randomized (Chapter 09, Section 6).
- The sample terminal session proved this concretely: editing one file, adding another, and deleting a third were all caught and correctly categorized by a single `fileguard check` run, with zero need to inspect file contents for anything "suspicious" by hand.
- This mini project is a direct, hands-on preview of Volume 3: a block's hash is this same fingerprinting idea applied to transactions, and a chain of blocks is this same "compare against last time" idea applied continuously and cumulatively.

---

## Exercises

### Easy

1. Build `fileguard` yourself, run `scan` against a small folder of your own files, and open the resulting `fileguard.json` in a text editor. Confirm every hash is exactly 64 hex characters long, and explain why (Chapter 09, Section 3).

2. Run `check` immediately after `scan` with no files changed, and confirm it reports "no changes detected." Then rename one file (without changing its contents) and run `check` again. Explain, in 3-4 sentences, why a rename is reported as one addition plus one deletion rather than as a "rename," referencing Section 7's diff logic.

3. Modify `printDiff` (Section 8) to also print a final summary line, such as `3 added, 1 changed, 2 deleted`, after listing individual files.

### Medium

4. Section 5 noted that `os.ReadFile` loads an entire file into memory at once. Rewrite `buildManifest`'s hashing step to use `sha256.New()` combined with `io.Copy` (streaming the file through the hasher in fixed-size chunks instead of loading it all at once), and confirm your rewritten version produces identical hashes to the original for the same files. Explain in 150-200 words why this version scales better for very large files.

5. Add a `-ignore` flag to `fileguard` that accepts a comma-separated list of directory names to skip entirely during `buildManifest`'s walk (for example, `-ignore=.git,node_modules`). Test it against a folder containing at least one directory you deliberately want excluded.

6. Section 10 mentioned detecting renamed files by noticing a deleted path's hash reappearing under a new path. Implement this: extend `DiffResult` with a `Renamed map[string]string` field (old path to new path), and update `diffManifests` to detect the case where a hash present in `old.Files` under one path is also present in `newM.Files` under a *different* path, removing that pair from `Added`/`Deleted` and adding it to `Renamed` instead.

### Hard

7. Design and implement a `.fileguardignore` file format (one ignore pattern per line, similar in spirit to `.gitignore`), loaded automatically by `buildManifest` if present at the scanned root. Support at minimum plain directory-name and file-name matches (glob patterns are a reasonable stretch goal, using Go's `path/filepath.Match`). Write a short explanation (150-250 words) of why silently skipping files during a security-relevant scan is itself a decision worth logging or reporting, not just implementing silently.

8. `fileguard` currently trusts its own `fileguard.json` manifest completely — if an attacker who tampered with a file also has write access to the manifest, they could update the manifest's stored hash to match their tampered file, and `check` would report no changes at all. Using Chapter 12 and Chapter 13's `crypto.Sign` and `crypto.Verify`, design (and implement, if you like) a `fileguard scan --sign` mode that signs the completed manifest with a private key, and a `fileguard check --verify` mode that refuses to trust a manifest whose signature does not verify. Explain, in 250-350 words, exactly what class of attack this addition defends against that a plain, unsigned manifest does not.

9. Research how real version-control systems (Git specifically) use content hashing as the foundation of their entire object model, not just for change detection but for storage itself (a "blob" in Git is essentially content addressed by its own hash). Write a comparison (300-400 words) of `fileguard`'s manifest-based approach against Git's content-addressed object store, and explain one specific advantage Git's approach has that a simple path-to-hash manifest, like this chapter's, does not.
