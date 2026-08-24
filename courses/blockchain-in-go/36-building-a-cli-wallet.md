# Chapter 36: Building a CLI Wallet

Every previous chapter in this course has been exercised through Go function calls in tests and demos. This chapter builds the first tool a real, non-programmer user would actually touch: a command-line wallet that generates keys, checks balances, and sends coins.

## Table of Contents

1. [What Real Users Actually Touch](#1-what-real-users-actually-touch)
2. [Designing the gochain-wallet Command](#2-designing-the-gochain-wallet-command)
3. [A Minimal Local Keystore](#3-a-minimal-local-keystore)
4. [Subcommand: new](#4-subcommand-new)
5. [Subcommand: balance](#5-subcommand-balance)
6. [Subcommand: send](#6-subcommand-send)
7. [Wiring the Subcommands Together in main.go](#7-wiring-the-subcommands-together-in-maingo)
8. [Trying It Out End to End](#8-trying-it-out-end-to-end)
9. [What's Missing (On Purpose, For Now)](#9-whats-missing-on-purpose-for-now)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What Real Users Actually Touch

Every chapter since Volume 2 has built GoChain from the inside: types, methods, tests, benchmarks. None of that is anything a normal person would ever type or run directly — nobody using a real cryptocurrency wallet writes Go code to call `wallet.New()`. What they actually experience is a small number of simple, human commands: "give me a new address," "how much do I have," "send some to this address." Everything underneath — key generation, UTXO scanning, transaction signing — is supposed to disappear behind those three requests.

This chapter builds exactly that surface: a `gochain-wallet` command-line tool with three subcommands, `new`, `balance`, and `send`. It is deliberately not fancy. Go's standard `flag` package (already used for the `gochain inspect` tool in Chapter 21) is enough for now; Chapter 74 in Volume 10 replaces every ad-hoc CLI in this course with a single, polished tool built on the Cobra framework. The goal here is narrower: prove that everything built across Volumes 2 through 5 — hashing, signatures, addresses, transactions, the mempool, fees — actually adds up to something a user can operate without reading a single line of the underlying Go code.

```
                    gochain-wallet
                          |
        +-----------------+-----------------+
        |                 |                 |
       new              balance            send
  (generate a       (scan the UTXO    (build, sign, and
   new key pair      set for an        submit a transfer
   and address)       address's        to the mempool)
                       balance)
```

---

## 2. Designing the gochain-wallet Command

The tool is invoked as `gochain-wallet <subcommand> [flags]`, following the same convention as `git <subcommand>` or `go <subcommand>` — a pattern Go developers already know instinctively. Each subcommand gets its own `flag.FlagSet`, so `-address` means something different (and takes different values) under `balance` than it might under `send`.

```
$ gochain-wallet new
Your new address: 1J9k...Qz7f
(private key saved to ./wallets/1J9k...Qz7f.json — keep this file safe!)

$ gochain-wallet balance -address 1J9k...Qz7f
Balance: 50 gochips

$ gochain-wallet send -from 1J9k...Qz7f -to 1Bv1...k7c -amount 20
Transaction abcd1234... submitted to the mempool.
```

Every subcommand ultimately needs the same two pieces of shared state: a local `*core.Blockchain` to read from, and a `*core.Mempool` to submit new transactions into. For this chapter, both are kept as simple, local, in-memory-per-run values — exactly matching the chapter's own scope ("submitting locally, for now"). Chapter 37's major project shows a complete run using this exact wallet against a shared chain and mempool; real networking, where `send` actually reaches other nodes, doesn't arrive until Volume 7.

---

## 3. A Minimal Local Keystore

Here's an immediate wrinkle: `new` and `send` are two *separate invocations* of the program, potentially minutes or days apart. Whatever private key `new` generates has to survive between those two runs somehow, or `send` would have no way to sign anything. Volume 6 builds GoChain's real answer to this — encrypted wallet files, HD key derivation, seed phrases — but none of that exists yet. This chapter needs something simple enough to write today, clearly labeled as a placeholder.

We'll store just enough to reconstruct a key pair: the private scalar (`D`) from the wallet's ECDSA key, alongside the curve it uses. Given `D` and the curve, the public key can always be recomputed — you never need to separately store the public half.

```go
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
)

// keystoreFile is a deliberately minimal, UNENCRYPTED on-disk format for
// this chapter's CLI wallet only. It exists purely so `new` and `send`
// can be run as separate program invocations and still share a key.
// Volume 6 replaces this entirely with encrypted, HD wallet files —
// never use this format for anything holding real value.
type keystoreFile struct {
	D []byte `json:"d"` // the private key's scalar value
}

func keystorePath(address string) string {
	return filepath.Join("wallets", address+".json")
}

// saveKey writes just enough of a private key to disk to reconstruct it
// later: the private scalar D. The public key and address are always
// deterministically re-derivable from D plus the curve, so there's
// nothing else that needs saving.
func saveKey(address string, priv *ecdsa.PrivateKey) error {
	if err := os.MkdirAll("wallets", 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(keystoreFile{D: priv.D.Bytes()})
	if err != nil {
		return err
	}
	return os.WriteFile(keystorePath(address), data, 0o600)
}

// loadKey reconstructs a private key from disk by reading back D and
// recomputing the matching public key point on the same curve GoChain's
// crypto package uses (established in Chapter 13).
func loadKey(address string) (*ecdsa.PrivateKey, error) {
	raw, err := os.ReadFile(keystorePath(address))
	if err != nil {
		return nil, fmt.Errorf("no local wallet found for address %s: %w", address, err)
	}

	var kf keystoreFile
	if err := json.Unmarshal(raw, &kf); err != nil {
		return nil, err
	}

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256() // the same curve chosen in Chapter 13
	priv.D = new(big.Int).SetBytes(kf.D)
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(kf.D)

	return priv, nil
}
```

`keystoreFile` is the on-disk shape — just the private scalar `D`, nothing else. `saveKey` creates a local `wallets/` directory and writes one JSON file per address, named after the address itself so `loadKey` can find it again later. `loadKey` does the reverse: it reads `D` back, builds a fresh `ecdsa.PrivateKey`, sets its curve to the same `elliptic.P256()` curve GoChain's crypto package settled on back in Chapter 13, and recomputes the public key's `X`/`Y` coordinates from `D` using `ScalarBaseMult` — exactly the same math that produced them the first time. The file permissions (`0o600`, `0o700`) are a nod toward "don't make this world-readable," but the loud comment at the top is the real safeguard: this format is explicitly unencrypted and explicitly temporary.

---

## 4. Subcommand: new

Generating a wallet reuses `wallet.New()` — the constructor this course's shared contract guarantees exists — and then persists its key using the keystore from Section 3:

```go
func runNew(args []string) {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	fs.Parse(args)

	w := wallet.New() // generates a fresh ECDSA key pair and derives an address (Chapters 13-14)

	if err := saveKey(w.Address(), w.KeyPair.PrivateKey); err != nil {
		log.Fatalf("failed to save the new wallet: %v", err)
	}

	fmt.Println("Your new address:", w.Address())
	fmt.Printf("(private key saved to %s — keep this file safe!)\n", keystorePath(w.Address()))
}
```

`runNew` takes a `flag.FlagSet` even though `new` has no flags of its own yet, purely for consistency with the other two subcommands and to leave room for future options (like `-outdir`). `wallet.New()` does all the cryptographic heavy lifting — generating a private/public key pair and deriving a checksummed Base58 address from it, exactly as built in Chapters 13 and 14. `saveKey` immediately persists the private key so a later `send` invocation can recover it, and the function finishes by printing the one thing a user actually needs to remember: their new address.

---

## 5. Subcommand: balance

Checking a balance means scanning the UTXO set for every unspent output belonging to a given address — the operation the course's own vision `main.go` (see the table of contents introduction) calls `chain.BalanceOf(address)`:

```go
func runBalance(bc *core.Blockchain, args []string) {
	fs := flag.NewFlagSet("balance", flag.ExitOnError)
	address := fs.String("address", "", "address to check the balance of")
	fs.Parse(args)

	if *address == "" {
		log.Fatal("balance: -address is required")
	}

	balance := bc.BalanceOf(*address)
	fmt.Printf("Balance: %d gochips\n", balance)
}
```

Under the hood, `BalanceOf` does exactly what Chapter 30 described conceptually: it derives the public-key hash encoded in `*address`, walks every unspent output tracked by the chain's UTXO set (built and kept current since Chapter 32), and sums the `Value` of every output whose `PubKeyHash` matches. None of that complexity leaks into this subcommand — `runBalance` just parses the `-address` flag, calls the one method that does the real work, and prints the result. This is exactly the point of building `core.Blockchain` and `core.UTXOSet` as clean, well-named types across the last several chapters: a CLI command that uses them correctly ends up almost boringly simple.

---

## 6. Subcommand: send

Sending is where every piece from this volume comes together at once: loading a key, building a transaction against the current UTXO set, and submitting it to the mempool.

```go
func runSend(bc *core.Blockchain, mempool *core.Mempool, utxoSet *core.UTXOSet, args []string) {
	fs := flag.NewFlagSet("send", flag.ExitOnError)
	from := fs.String("from", "", "sender's address (must have a local wallet file)")
	to := fs.String("to", "", "recipient's address")
	amount := fs.Int64("amount", 0, "amount to send, in gochips")
	fs.Parse(args)

	if *from == "" || *to == "" || *amount <= 0 {
		log.Fatal("send: -from, -to, and a positive -amount are all required")
	}

	priv, err := loadKey(*from)
	if err != nil {
		log.Fatalf("send: %v", err)
	}
	w := &wallet.Wallet{KeyPair: &crypto.KeyPair{PrivateKey: priv}}

	// NewTransaction (Chapter 32) selects enough UTXOs to cover the
	// amount, creates a change output back to the sender if needed, and
	// signs every input using the wallet's private key before returning.
	tx, err := core.NewTransaction(w, *to, *amount, utxoSet)
	if err != nil {
		log.Fatalf("send: could not build transaction: %v", err)
	}

	// "Submitting," for now, means handing it to this same process's
	// mempool. Real broadcasting to other nodes arrives in Volume 7.
	if err := mempool.Add(tx); err != nil {
		log.Fatalf("send: transaction rejected by mempool: %v", err)
	}

	fmt.Printf("Transaction %x submitted to the mempool.\n", tx.ID)
}
```

`runSend` first parses and validates its three flags — `-from`, `-to`, `-amount` — refusing to proceed if any of them looks obviously wrong (a missing address, a zero or negative amount). It then calls `loadKey` from Section 3 to recover the sender's private key from disk, and wraps it into a minimal `wallet.Wallet` value just shaped enough to satisfy `core.NewTransaction`'s signature. The call to `core.NewTransaction` does the real work this entire volume has been building toward: it consults `utxoSet` to find enough of the sender's unspent outputs to cover `*amount`, builds a change output back to the sender if the selected UTXOs add up to more than `*amount`, and signs every input with the sender's private key before returning a fully-formed, ready-to-broadcast `*core.Transaction`. The final step — `mempool.Add(tx)` — is exactly the Chapter 34 mempool doing its job: accepting the transaction if nothing else pending conflicts with it, or rejecting it (for instance, as a double-spend) if it does.

---

## 7. Wiring the Subcommands Together in main.go

Here is the complete, runnable entry point tying the three subcommands together:

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/you/gochain/core"
	"github.com/you/gochain/crypto"
	"github.com/you/gochain/wallet"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: gochain-wallet <new|balance|send> [flags]")
		os.Exit(1)
	}

	// bc, mempool, and utxoSet represent this run's local view of the
	// chain. In Chapter 37's demo they're shared across an entire
	// end-to-end scenario in one process; here, imagine them backed by
	// whatever local chain state a real node would already have loaded.
	bc := core.NewBlockchain()
	mempool := core.NewMempool()
	utxoSet := core.NewUTXOSet(bc)
	utxoSet.Reindex() // rebuild the UTXO set from the chain, from scratch

	switch os.Args[1] {
	case "new":
		runNew(os.Args[2:])
	case "balance":
		runBalance(bc, os.Args[2:])
	case "send":
		runSend(bc, mempool, utxoSet, os.Args[2:])
	default:
		fmt.Printf("unknown subcommand %q\n", os.Args[1])
		fmt.Println("usage: gochain-wallet <new|balance|send> [flags]")
		os.Exit(1)
	}
}
```

The dispatch logic is intentionally plain: look at `os.Args[1]` (the first word after the program name) and route to `runNew`, `runBalance`, or `runSend` accordingly, handing each one whatever slice of the remaining arguments (`os.Args[2:]`) it needs to parse its own flags. `core.NewBlockchain()`, `core.NewMempool()`, and `core.NewUTXOSet(bc)` set up this run's local view of chain state; `utxoSet.Reindex()` populates it by scanning the chain that already exists, the same rebuild-from-scratch operation Chapter 37 leans on heavily for verifying balances independently.

---

## 8. Trying It Out End to End

With the pieces in place, here's what an actual session looks like (assuming a block has already rewarded the first address with 50 gochips, exactly as Chapter 37's demo sets up):

```bash
$ go build -o gochain-wallet .

$ ./gochain-wallet new
Your new address: 1AbCdEfGhIjKlMnOpQrStUvWxYz1234567
(private key saved to wallets/1AbCdEfGhIjKlMnOpQrStUvWxYz1234567.json — keep this file safe!)

$ ./gochain-wallet new
Your new address: 1ZzYyXxWwVvUuTtSsRrQqPpOoNnMmLlKkJj
(private key saved to wallets/1ZzYyXxWwVvUuTtSsRrQqPpOoNnMmLlKkJj.json — keep this file safe!)

$ ./gochain-wallet balance -address 1AbCdEfGhIjKlMnOpQrStUvWxYz1234567
Balance: 50 gochips

$ ./gochain-wallet balance -address 1ZzYyXxWwVvUuTtSsRrQqPpOoNnMmLlKkJj
Balance: 0 gochips

$ ./gochain-wallet send -from 1AbCdEfGhIjKlMnOpQrStUvWxYz1234567 -to 1ZzYyXxWwVvUuTtSsRrQqPpOoNnMmLlKkJj -amount 20
Transaction 9f86d081884c7d659a2feaa0c55ad015... submitted to the mempool.
```

Two calls to `new` generate two independent local wallets, each with its own saved keystore file. `balance` confirms the first address already has funds (from a previously mined block) and the second starts at zero. `send` then builds, signs, and submits a 20-gochip transfer from the first address to the second — at this point the transaction is sitting in the mempool, not yet reflected in either address's balance, since nothing has mined it into a block yet. That's precisely the gap Chapter 37 closes: mining the pending transaction and re-checking both balances afterward to confirm they moved exactly as expected.

---

## 9. What's Missing (On Purpose, For Now)

Being upfront about scope keeps expectations calibrated for what comes next:

- **No real networking.** `send` submits to this process's own local mempool. A separate `gochain-wallet` process — let alone a wallet running on someone else's computer — has no way to see it yet. Volume 7 builds the P2P layer that makes "submit" actually mean "broadcast to the network."
- **No encryption, no seed phrases, no multiple addresses per wallet.** Section 3's keystore is a stopgap. Volume 6 replaces it with BIP-39 seed phrases, BIP-32/44 hierarchical key derivation, and password-encrypted wallet files — the same techniques MetaMask and hardware wallets use.
- **No mining command.** This CLI can generate keys, check balances, and submit transactions, but it doesn't mine blocks itself. In Chapter 37's demo, mining is driven directly rather than through this CLI; a real `gochain node` command that mines continuously is a natural extension you're encouraged to try in the exercises.
- **No persistence for the blockchain or mempool across runs.** Every invocation of this CLI in this chapter starts a fresh, empty `core.NewBlockchain()`. That's fine for illustrating the three subcommands independently, but it means, as written, `balance` and `send` in Section 8 wouldn't actually see each other's effects across separate process runs without a shared, persisted chain — exactly why Chapter 37 runs its full demo inside one continuous program rather than as separate CLI invocations.

---

## Summary

- The `gochain-wallet` CLI is the first tool in this course a real end user would actually operate, hiding key generation, UTXO scanning, and transaction signing behind three subcommands: `new`, `balance`, `send`.
- Each subcommand gets its own `flag.FlagSet`, following the same "program subcommand [flags]" convention as `git` and `go`.
- A minimal, explicitly unencrypted local keystore (just the private scalar `D`, plus the known curve) lets `new` and `send` share a key across separate process runs, standing in for the real, encrypted wallet files Volume 6 builds.
- `new` calls `wallet.New()` and persists the resulting key; `balance` calls `core.Blockchain.BalanceOf(address)`, which scans the UTXO set for matching outputs; `send` loads a saved key, calls `core.NewTransaction` to build and sign a transfer, and submits it via `mempool.Add`.
- `send` in this chapter only reaches this process's own local mempool — real broadcasting to other nodes is Volume 7's job.
- Trying the three subcommands together demonstrates the entire transaction pipeline from Volumes 2 through 5 operating behind a simple, human-usable interface.
- This CLI is intentionally a stepping stone: Chapter 42 adds a full HD wallet CLI, and Chapter 74 rebuilds the whole CLI surface on Cobra.

---

## Exercises

### Easy

1. Add a fourth subcommand, `address`, that takes no flags and simply lists every address currently saved under `./wallets/` (by reading the directory's filenames), so a user can rediscover which wallets they've already generated.
2. Modify `runBalance` to print a friendlier error message (instead of relying on `core.Blockchain.BalanceOf`'s default behavior) when the given address string fails to decode as a valid GoChain address at all — distinguishing "this address has zero balance" from "this isn't a valid address."
3. Add a `-help` flag behavior (or reuse each `FlagSet`'s default `-h`) that prints a one-line usage example for each subcommand, and verify `gochain-wallet send -h` prints something a first-time user could actually follow.

### Medium

4. `runSend` currently fails hard (`log.Fatalf`) if `loadKey` can't find a keystore file for `-from`. Change it to return a proper Go `error` up through `main` instead of calling `log.Fatalf` deep inside a helper, and have `main` be the only place that ever calls `os.Exit`, following the Go convention of keeping library-style code testable and free of direct process-exit calls.
5. Implement the missing piece hinted at in Section 9: a fourth subcommand, `mine`, that takes no address-specific flags beyond `-miner <address>`, pulls everything pending from the shared mempool, builds a coinbase transaction rewarding `-miner`, and calls `bc.MineBlock` — wiring in the fee-based selection from Chapter 35 rather than just passing every pending transaction through unfiltered.
6. Right now, every invocation of `gochain-wallet` starts a brand-new, empty in-memory blockchain (Section 9's last bullet). Add a simple persistence layer (even a single gob-encoded file is fine, foreshadowing Chapter 20's approach) so that `balance` and `send` invoked as separate process runs actually see a consistent, shared chain and mempool state.

### Hard

7. Extend the keystore from Section 3 to support multiple wallets under one arbitrary "profile name" rather than being indexed purely by address, and add a `-wallet <name>` flag usable across all three subcommands, so a user could maintain, say, a "savings" and a "spending" wallet side by side without needing to memorize raw addresses.
8. Add input validation to `send` that rejects a transaction attempt where `-amount` would require selecting an unreasonably large number of UTXOs (say, more than 100) to satisfy, printing a helpful suggestion to "consolidate" funds first — and implement a `consolidate` subcommand that builds a transaction spending many small UTXOs into one larger one sent back to the same address.
9. Turn this chapter's keystore format into a genuinely safer stopgap without going all the way to Volume 6's full HD wallet system: prompt for a passphrase on `new` and `send` (reading from stdin without echoing it to the terminal), and use it to encrypt the saved `D` value with a standard authenticated encryption scheme (`crypto/aes` in GCM mode) before writing it to disk, decrypting it again in `loadKey`. Write a test proving a wrong passphrase fails to decrypt rather than silently producing a garbage key.
