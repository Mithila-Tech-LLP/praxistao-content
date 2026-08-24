# Chapter 74: Building a Polished CLI with Cobra

Across this course, GoChain has accumulated a CLI one chapter at a time: `gochain inspect` in Chapter 21, a wallet tool in Chapter 36, an HD wallet CLI in Chapter 42, none of them sharing a common command structure, help system, or flag convention. Each was built with the standard library's `flag` package, which is perfectly fine for one small standalone tool, but starts to creak the moment you want several related tools to feel like a single, coherent product. This chapter replaces every one of those ad-hoc entry points with a single `gochain` command, built on Cobra — the same library behind `kubectl`, `docker`, `hugo`, and `go` itself — so that `gochain node start`, `gochain wallet new`, `gochain tx send`, and `gochain chain inspect` all live under one root, with consistent help text, consistent flags, and room to grow.

## Table of Contents

1. [Why Move Beyond `flag`](#1-why-move-beyond-flag)
2. [Installing Cobra and Project Layout](#2-installing-cobra-and-project-layout)
3. [The Root Command](#3-the-root-command)
4. [`gochain node start`](#4-gochain-node-start)
5. [`gochain wallet new`](#5-gochain-wallet-new)
6. [`gochain tx send`](#6-gochain-tx-send)
7. [`gochain chain inspect`](#7-gochain-chain-inspect)
8. [Wiring It All Together in `main.go`](#8-wiring-it-all-together-in-maingo)
9. [Trying It Out](#9-trying-it-out)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why Move Beyond `flag`

Picture a toolbox where every tool has its own, slightly different handle: one wrench you grip overhand, another underhand, a third with a trigger nobody explained. That's what GoChain's CLI has become — `gochain-inspect --verify`, a separate `gochain-wallet new` binary, a separate `gochain-wallet send` binary — each perfectly usable alone, each requiring you to relearn its particular conventions. A polished CLI feels like a single, well-designed toolbox instead: one entry point (`gochain`), a consistent way to discover what it can do (`gochain help`, `gochain wallet --help`), and a predictable grammar — `gochain <noun> <verb> [flags]` — you can guess your way through even for a subcommand you've never used.

```
BEFORE (this course, Chapters 21, 36, 42)          AFTER (this chapter)

  gochain-inspect --verify                          gochain chain inspect --verify
  gochain-wallet new                                gochain wallet new
  gochain-wallet send --to ... --amount ...          gochain tx send --to ... --amount ...
  (separate binary, separate main())                 gochain node start --port 8080

  three+ binaries, three+ flag conventions           one binary, one command tree,
                                                      one consistent --help everywhere
```

Cobra is a good fit here for a concrete reason beyond consistency: it generates working `--help` text, usage strings, and even shell autocompletion *for free*, directly from the command tree you declare — the same command tree becomes the documentation, rather than documentation you write and maintain separately and which inevitably drifts out of sync with the code.

---

## 2. Installing Cobra and Project Layout

```bash
go get github.com/spf13/cobra@latest
```

GoChain's CLI code moves into its own directory, separate from the library packages (`core`, `consensus`, `wallet`, and so on) it calls into — the same "keep `cmd/` separate from library code" principle Chapter 06 established for the module as a whole.

```
gochain/cmd/gochain/
├── main.go            -- entry point; just calls cmd.Execute()
├── root.go             -- the root `gochain` command and shared flags
├── node.go              -- `gochain node` and its subcommands
├── wallet.go             -- `gochain wallet` and its subcommands
├── tx.go                  -- `gochain tx` and its subcommands
└── chain.go                -- `gochain chain` and its subcommands
```

Every file above defines one or more `*cobra.Command` values and registers them onto a parent via `AddCommand` — the entire CLI is, structurally, just a tree of commands, exactly the way a filesystem is a tree of directories.

```
                        gochain (root)
                    /       |       |       \
               node      wallet    tx      chain
                |          |        |        |
              start       new     send    inspect
```

---

## 3. The Root Command

The root command is the tree's trunk: it owns global flags every subcommand can see (like `--data-dir`, where GoChain stores its blockchain and wallet files) and prints a short banner when `gochain` is run with no subcommand at all.

```go
// gochain/cmd/gochain/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// dataDir is a package-level variable bound to the --data-dir persistent
// flag below. "Persistent" in Cobra's vocabulary means every subcommand
// nested under rootCmd inherits this flag automatically -- gochain node
// start and gochain wallet new both see the same --data-dir without
// redeclaring it themselves.
var dataDir string

// rootCmd is the trunk of the whole command tree. Running "gochain" with
// no subcommand falls through to Run, which just prints a short banner --
// real work only happens inside the subcommands added in Section 4 onward.
var rootCmd = &cobra.Command{
	Use:   "gochain",
	Short: "GoChain: a blockchain node, wallet, and toolkit",
	Long: `GoChain is a from-scratch blockchain implementation in Go.

This single binary bundles everything you need to run a node, manage a
wallet, send transactions, and inspect the chain -- use "gochain <command>
--help" to see the flags and behavior of any specific command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoChain -- run 'gochain --help' to see available commands")
	},
}

// Execute is the single function main.go calls. Cobra's own Execute()
// parses os.Args, walks the command tree to find the right subcommand,
// and calls its Run function -- main.go itself never touches os.Args
// or flag parsing directly, which is the whole point of handing that
// job to Cobra.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Cobra already printed a usage message for the failing command;
		// os.Exit(1) just gives shell scripts a proper non-zero status
		// to check, the same convention every well-behaved CLI follows.
		os.Exit(1)
	}
}

func init() {
	// PersistentFlags (as opposed to Flags) are inherited by every
	// subcommand nested anywhere under rootCmd, no matter how deep --
	// gochain node start and gochain wallet new both see --data-dir.
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "./data",
		"directory where GoChain stores blockchain and wallet data")
}
```

`Use`, `Short`, and `Long` are not decoration — Cobra reads them directly to generate `gochain --help`'s output, so writing a clear `Long` description here is equivalent to writing a man page, except it can never drift out of sync with the actual command, because it lives in the same file as the command's behavior.

---

## 4. `gochain node start`

`node` is a **parent command with no `Run` of its own** — running bare `gochain node` should just show its own subcommands' help text, exactly like `git remote` (with no further arguments) lists what you can do with remotes rather than doing anything itself.

```go
// gochain/cmd/gochain/node.go
package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/you/gochain/api"
	"github.com/you/gochain/consensus"
	"github.com/you/gochain/core"
	"github.com/you/gochain/storage"
)

// nodeCmd has no Run function -- Cobra automatically shows its
// subcommands' help when it's invoked directly (`gochain node`), which
// is exactly the behavior a purely-organizational parent command wants.
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Run and manage a GoChain node",
}

var (
	nodeAddr    string
	nodeMinerAddr string
)

// nodeStartCmd is where Chapter 70's api.Server actually gets started --
// everything before this chapter built the pieces; this is where they
// get wired together into a single running process.
var nodeStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a GoChain node and its HTTP API",
	Long: `Start opens (or creates) the blockchain in --data-dir, wires up the
mempool and UTXO set, and serves the JSON-RPC/REST API from Chapter 70 on
the address given by --addr. It blocks until the process is stopped.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chain, err := core.OpenBlockchain(dataDir, nodeMinerAddr)
		if err != nil {
			return err
		}
		defer chain.Close()

		mempool := core.NewMempool()
		utxoSet, err := storage.NewUTXOSet(chain)
		if err != nil {
			return err
		}
		_ = consensus.NewProofOfWork // engine selection wired in at MineBlock call sites

		server := api.NewServer(chain, mempool, utxoSet)
		router := api.NewRouter(server)

		log.Printf("gochain node listening on %s (data dir: %s)", nodeAddr, dataDir)
		return http.ListenAndServe(nodeAddr, router)
	},
}

func init() {
	nodeStartCmd.Flags().StringVar(&nodeAddr, "addr", ":8080",
		"address the HTTP API listens on")
	nodeStartCmd.Flags().StringVar(&nodeMinerAddr, "miner-address", "",
		"address that receives coinbase rewards when this node mines a block")
	nodeStartCmd.MarkFlagRequired("miner-address")

	nodeCmd.AddCommand(nodeStartCmd)
	rootCmd.AddCommand(nodeCmd)
}
```

Note `RunE` instead of `Run` here — Cobra supports both, but `RunE` (which returns an `error`) is the better choice whenever a command can genuinely fail, like `core.OpenBlockchain` returning an error for a corrupted data directory. Returning the error from `RunE` lets Cobra print it consistently and set a non-zero exit code, rather than every command hand-rolling its own `log.Fatal` call. `MarkFlagRequired` is a small but meaningful Cobra feature: it makes `--miner-address` mandatory and produces a clear error ("required flag(s) \"miner-address\" not set") if it's missing, instead of the node silently starting with an empty reward address.

---

## 5. `gochain wallet new`

```go
// gochain/cmd/gochain/wallet.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/you/gochain/wallet"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Create and manage GoChain wallets",
}

var walletPassword string

// walletNewCmd wraps wallet.New (Volume 6) plus the encrypted-storage
// machinery from Chapter 40 -- generating a fresh HD wallet, saving it
// encrypted under --data-dir, and printing the address and (once,
// clearly labeled) the recovery seed phrase.
var walletNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a new wallet and print its address",
	Long: `New generates a fresh HD wallet (Chapter 38), encrypts it with
--password before saving it under --data-dir/wallet.dat, and prints the
wallet's first address plus its recovery seed phrase. The seed phrase is
shown exactly once -- write it down, because GoChain never displays it
again.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if walletPassword == "" {
			return fmt.Errorf("--password is required to encrypt the new wallet file")
		}

		w, mnemonic, err := wallet.NewHD()
		if err != nil {
			return err
		}
		if err := w.SaveEncrypted(dataDir+"/wallet.dat", walletPassword); err != nil {
			return err
		}

		fmt.Println("New wallet created.")
		fmt.Println("Address:      ", w.Address())
		fmt.Println("Seed phrase:  ", mnemonic)
		fmt.Println()
		fmt.Println("Write the seed phrase down somewhere safe -- it is the only way")
		fmt.Println("to recover this wallet, and GoChain will never show it again.")
		return nil
	},
}

func init() {
	walletNewCmd.Flags().StringVar(&walletPassword, "password", "",
		"password used to encrypt the new wallet file at rest")
	walletCmd.AddCommand(walletNewCmd)
	rootCmd.AddCommand(walletCmd)
}
```

`gochain wallet new --password ...` is a direct, single-command replacement for Chapter 42's standalone `gochain-wallet new`, with one meaningful improvement: because it's a subcommand of the same root as `gochain node start`, both commands automatically agree on `--data-dir` through the persistent flag from Section 3 — a wallet created here is saved to exactly the directory a node started later will look for it in, with no risk of the two tools disagreeing about a path.

---

## 6. `gochain tx send`

```go
// gochain/cmd/gochain/tx.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/you/gochain/wallet"
)

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Build, sign, and submit transactions",
}

var (
	txTo       string
	txAmount   int64
	txPassword string
	txRPCAddr  string
)

// txSendCmd loads the wallet from --data-dir, builds and signs a
// transaction using the exact wallet + core.NewTransaction machinery
// from Volume 5 and 6, and submits it to a running node's REST API
// (Chapter 70's POST /transactions) rather than adding it to a local
// mempool directly -- this command is a *client* of a node, the same
// way any external wallet app would be.
var txSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send gochips to another address",
	Long: `Send loads the wallet at --data-dir/wallet.dat (decrypting it with
--password), builds and signs a transaction paying --amount gochips to
--to, and submits it to the node at --rpc-addr via POST /transactions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if txTo == "" {
			return fmt.Errorf("--to is required")
		}
		if txAmount <= 0 {
			return fmt.Errorf("--amount must be a positive number of gochips")
		}

		w, err := wallet.LoadEncrypted(dataDir+"/wallet.dat", txPassword)
		if err != nil {
			return fmt.Errorf("failed to load wallet: %w", err)
		}

		tx, err := w.BuildAndSignTransaction(txTo, txAmount, txRPCAddr)
		if err != nil {
			return fmt.Errorf("failed to build transaction: %w", err)
		}

		txID, err := wallet.SubmitTransaction(txRPCAddr, tx)
		if err != nil {
			return fmt.Errorf("node rejected transaction: %w", err)
		}

		fmt.Println("Transaction submitted:", txID)
		fmt.Println("Track it at:", txRPCAddr+"/explorer/transactions/"+txID)
		return nil
	},
}

func init() {
	txSendCmd.Flags().StringVar(&txTo, "to", "", "recipient address (required)")
	txSendCmd.Flags().Int64Var(&txAmount, "amount", 0, "amount in gochips (required)")
	txSendCmd.Flags().StringVar(&txPassword, "password", "", "wallet password")
	txSendCmd.Flags().StringVar(&txRPCAddr, "rpc-addr", "http://localhost:8080",
		"address of the GoChain node to submit the transaction to")

	txCmd.AddCommand(txSendCmd)
	rootCmd.AddCommand(txCmd)
}
```

`gochain tx send` deliberately talks to a node over HTTP (via `wallet.SubmitTransaction`, a thin wrapper around Chapter 70's `POST /transactions`) rather than importing `core.Mempool` directly, even though both live in the same Go module. This mirrors how a real user would run things in practice: the wallet and the node are very often two separate processes — possibly on two separate machines — and building the CLI to talk over the same HTTP boundary every external client uses means there is exactly one code path for "submit a transaction," tested once, used everywhere.

---

## 7. `gochain chain inspect`

```go
// gochain/cmd/gochain/chain.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/you/gochain/core"
)

var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Inspect the local blockchain",
}

var chainVerify bool

// chainInspectCmd is Chapter 21's gochain-inspect, unchanged in behavior,
// just re-homed under the unified command tree with the same --verify
// flag it always had.
var chainInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Print every block in the chain, optionally verifying links",
	Long: `Inspect opens the chain at --data-dir and prints each block's height,
timestamp, transaction count, hash, and previous hash. With --verify, it
also walks the whole chain re-checking every hash link (Chapter 19),
refusing to report success if any block fails validation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chain, err := core.OpenBlockchain(dataDir, "")
		if err != nil {
			return err
		}
		defer chain.Close()

		iter := chain.Iterator()
		for {
			block, ok := iter.Next()
			if !ok {
				break
			}
			fmt.Printf("Height:      %d\n", block.Height)
			fmt.Printf("Timestamp:   %d\n", block.Timestamp)
			fmt.Printf("Tx count:    %d\n", len(block.Transactions))
			fmt.Printf("Hash:        %x\n", block.Hash)
			fmt.Printf("Prev hash:   %x\n", block.PrevBlockHash)
			fmt.Println("---")
		}

		if chainVerify {
			if err := core.ValidateChain(chain); err != nil {
				return fmt.Errorf("chain validation FAILED: %w", err)
			}
			fmt.Println("Chain verified OK: every block links correctly.")
		}
		return nil
	},
}

func init() {
	chainInspectCmd.Flags().BoolVar(&chainVerify, "verify", false,
		"walk the whole chain re-checking every hash link")
	chainCmd.AddCommand(chainInspectCmd)
	rootCmd.AddCommand(chainCmd)
}
```

This is the smallest change of the four commands in this chapter, deliberately: Chapter 21's logic was already correct, so the only thing worth changing is *where it lives* — as `gochain chain inspect` instead of a standalone `gochain-inspect` binary, sharing the same `--data-dir` convention as every other command in the tree.

---

## 8. Wiring It All Together in `main.go`

With every subcommand registering itself onto `rootCmd` inside its own file's `init()` function, `main.go` shrinks down to almost nothing — its only job is calling `cmd.Execute()`.

```go
// gochain/cmd/gochain/main.go
package main

import "github.com/you/gochain/cmd"

func main() {
	cmd.Execute()
}
```

This is the payoff of Cobra's design: every command file is self-registering via `init()`, so `main.go` never needs to know the full list of commands, flags, or subcommands that exist — adding a brand-new `gochain contract deploy` command in some later chapter would mean adding one new file to `cmd/gochain/`, with zero changes required here.

---

## 9. Trying It Out

```bash
$ go build -o gochain ./cmd/gochain

$ ./gochain --help
GoChain is a from-scratch blockchain implementation in Go.

This single binary bundles everything you need to run a node, manage a
wallet, send transactions, and inspect the chain -- use "gochain <command>
--help" to see the flags and behavior of any specific command.

Usage:
  gochain [command]

Available Commands:
  chain       Inspect the local blockchain
  node        Run and manage a GoChain node
  tx          Build, sign, and submit transactions
  wallet      Create and manage GoChain wallets

Flags:
      --data-dir string   directory where GoChain stores blockchain and wallet data (default "./data")
  -h, --help              help for gochain

$ ./gochain wallet new --password hunter2
New wallet created.
Address:       gochain1alice7f3e9...
Seed phrase:    hollow region civil ... (12 words)

Write the seed phrase down somewhere safe -- it is the only way
to recover this wallet, and GoChain will never show it again.

$ ./gochain node start --miner-address gochain1alice7f3e9... --addr :8080
2026/08/01 09:12:03 gochain node listening on :8080 (data dir: ./data)

# in a second terminal:
$ ./gochain tx send --to gochain1ben4d2a1... --amount 25 --password hunter2
Transaction submitted: 9f3a7c1e...
Track it at: http://localhost:8080/explorer/transactions/9f3a7c1e...

$ ./gochain chain inspect --verify
Height:      0
Timestamp:   1743500000
Tx count:    1
Hash:        00000000...
Prev hash:   00000000...
---
Chain verified OK: every block links correctly.
```

Every one of these commands shares the same `--data-dir` default, the same help conventions, and the same non-zero exit code on failure — a small thing individually, but the accumulation of small, consistent things is precisely what "polished" means for a CLI.

---

## Summary

- Cobra replaces GoChain's collection of separate, `flag`-based binaries with one `gochain` command and a tree of subcommands, matching the pattern used by tools like `kubectl`, `docker`, and `go` itself.
- The command tree is `gochain` (root) → `node`/`wallet`/`tx`/`chain` (parent commands, no `Run` of their own) → `start`/`new`/`send`/`inspect` (leaf commands that do the actual work).
- `PersistentFlags` on the root command (like `--data-dir`) are automatically inherited by every subcommand nested anywhere underneath it, so every command agrees on where wallet and chain data lives without repeating the flag.
- `RunE` (rather than `Run`) lets a command return a Go `error`, which Cobra reports consistently and turns into a non-zero exit code — better than each command hand-rolling its own `log.Fatal`.
- `gochain node start` wires together `core.OpenBlockchain`, `core.Mempool`, `storage.UTXOSet`, and `api.NewServer`/`api.NewRouter` from Chapter 70 into one running process.
- `gochain tx send` deliberately talks to a node over the same `POST /transactions` HTTP endpoint (Chapter 70) any external wallet would use, rather than importing `core.Mempool` directly, even though both live in the same module.
- Every command file self-registers onto `rootCmd` via its own `init()` function, so `main.go` shrinks to a single call to `cmd.Execute()` and never needs an exhaustive list of commands.
- `Use`, `Short`, and `Long` fields on each command double as both the command's behavior declaration and its `--help` text — documentation that cannot silently drift out of sync with the code, because both live in the same place.

---

## Exercises

### Easy

1. Add a `--json` flag to `gochain chain inspect` that, when set, prints each block as a JSON object (one per line) instead of the human-readable text format, useful for piping into `jq`.
2. Cobra supports a `version` subcommand almost for free via `rootCmd.Version`. Add a version string (even a hardcoded one like `"0.1.0"`) to `rootCmd` and confirm `gochain --version` and `gochain version` both work.
3. `gochain wallet new` currently requires `--password` via a manual `if walletPassword == ""` check inside `RunE`. Replace it with Cobra's `MarkFlagRequired`, matching the pattern `nodeStartCmd` already uses for `--miner-address`, and confirm the error message Cobra produces on its own.

### Medium

4. Add a `gochain wallet balance --address <addr>` command that calls a running node's `GET /balance/{address}` endpoint (Chapter 70) and prints the result, reusing the same `--rpc-addr` flag convention `tx send` established.
5. Add shell autocompletion: Cobra generates completion scripts for bash, zsh, and fish automatically via `rootCmd.GenBashCompletionFile` (and friends). Add a hidden `gochain completion <shell>` command (Cobra has a built-in helper for this) and document, in a short comment, how a user would install the generated script.
6. Currently, `gochain tx send`'s `--password` flag takes the password directly on the command line, which any other user on the same machine could see in the process list (`ps aux`). Add a `--password-stdin` flag that instead reads the password from standard input, and update the help text to recommend it over `--password` for anything beyond local testing.

### Hard

7. Refactor the `nodeCmd`/`walletCmd`/`txCmd`/`chainCmd` `init()`-based self-registration into an explicit `NewRootCommand()` constructor function that builds and returns the whole tree without relying on package-level `init()` side effects, and explain in a comment why this makes the CLI easier to unit test (hint: consider trying to run two independent command trees with different flag values inside the same test binary).
8. Add a `gochain config` command group (`gochain config set key value`, `gochain config get key`) that reads and writes a small YAML or TOML file under `--data-dir`, and change `nodeStartCmd`'s `--addr` and `--miner-address` flags to fall back to a value from this config file when not passed explicitly on the command line, following the common CLI convention "flag overrides config file overrides default."
9. Write an integration test (using Go's `os/exec` to actually invoke the compiled `gochain` binary as a subprocess) that runs `gochain wallet new`, extracts the printed address, starts `gochain node start` in the background, submits a transaction with `gochain tx send`, and asserts `gochain chain inspect --verify` eventually reports the transaction mined and the chain valid — a full command-line-driven version of Chapter 37's end-to-end demo.
