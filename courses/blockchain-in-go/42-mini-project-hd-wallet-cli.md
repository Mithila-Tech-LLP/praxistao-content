# Chapter 42: Mini Project — HD Wallet CLI

Chapter 36 built a first, deliberately minimal `gochain-wallet` CLI, with a keystore format its own documentation labeled "explicitly unencrypted and explicitly temporary." Chapters 38 through 41 spent this entire volume building the real replacement: seed phrases, hierarchical key derivation, and password-encrypted storage. This chapter cashes in every one of those pieces at once, rebuilding `gochain-wallet` from the ground up on top of them — a CLI a real user could actually trust with real value, at least within this course's scope.

## Table of Contents

1. [From Single-Key CLI to HD Wallet CLI](#1-from-single-key-cli-to-hd-wallet-cli)
2. [Designing the Five Subcommands](#2-designing-the-five-subcommands)
3. [Encrypting a Mnemonic Instead of a Single Key](#3-encrypting-a-mnemonic-instead-of-a-single-key)
4. [Reading Passwords Without Echoing Them](#4-reading-passwords-without-echoing-them)
5. [Implementing new and recover](#5-implementing-new-and-recover)
6. [Implementing addresses, balance, and send](#6-implementing-addresses-balance-and-send)
7. [Mini Project: gochain-wallet](#mini-project-gochain-wallet)
8. [What's Still Missing](#8-whats-still-missing)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. From Single-Key CLI to HD Wallet CLI

Chapter 36's keystore stored exactly one thing: a single private key's raw scalar, in plaintext, on disk. That was enough to demonstrate `new`, `balance`, and `send` working together, but it left three problems on the table, each one this volume already solved at the library level:

```
Chapter 36's gochain-wallet              This chapter's gochain-wallet

  ONE key pair per wallet file             ONE seed backs an ENTIRE tree
  -> new address needs a new file          -> "addresses" lists as many
                                               as you want, from one backup

  PLAINTEXT private key on disk            ENCRYPTED with a user password
  -> a stolen file is a stolen wallet      -> a stolen file alone is useless
                                               without the password too

  Backup = "keep this file safe forever"   Backup = "remember twelve words"
  -> lose the file, lose the coins         -> the SEED PHRASE is the real
                                               backup; the file is just a
                                               convenience cache of it
```

Every capability on the right side of that table already exists as a tested, working function in `gochain/wallet` — `NewMnemonic`, `SeedFromMnemonic`, `NewHDWallet`, `DeriveWallet`, and the scrypt/AES-GCM machinery behind `EncryptWalletFile`/`DecryptWalletFile`. This chapter's actual job is almost entirely *plumbing*: connecting those functions to five commands a real person would type.

---

## 2. Designing the Five Subcommands

```
                        gochain-wallet
                              |
        +---------+----------+----------+---------+
        |         |          |          |         |
       new     recover   addresses   balance    send
   (generate   (restore  (list        (check     (build,
    a fresh     from an   derived      one         sign, and
    seed        existing  addresses)   address's   submit a
    phrase)     phrase)                balance)    transfer)
```

- **`new`** — generates a fresh BIP-39 mnemonic (Chapter 39), asks the user to choose a password, saves the mnemonic encrypted to disk (Chapter 40's crypto, applied to a mnemonic instead of a single key — Section 3 shows exactly how), and prints the seed phrase once plus a handful of derived addresses.
- **`recover`** — the mirror image: takes an *existing* mnemonic (typed in by the user, standing in for "written down on paper months ago"), re-encrypts it under a (possibly new) password, and proves recovery works by printing the same addresses `new` would have shown originally.
- **`addresses`** — unlocks the saved wallet file and lists as many derived addresses as requested, by index, using `DeriveWallet`.
- **`balance`** — unlocks the wallet, derives one specific address by index, and checks its balance — functionally identical to Chapter 36's `balance`, just backed by an HD-derived wallet instead of a single stored key.
- **`send`** — unlocks the wallet, derives the sending address by index, builds and signs a transaction with `core.NewTransaction` and `Transaction.Sign` (Chapters 32-33), and submits it to a mempool — again, functionally identical to Chapter 36's `send`, just sourcing its private key from HD derivation instead of a flat file.

Every command shares one on-disk file, `gochain-wallet.dat` by default (overridable with `-keystore`), which never stores a mnemonic in the clear at rest.

---

## 3. Encrypting a Mnemonic Instead of a Single Key

Chapter 40's `EncryptWalletFile`/`DecryptWalletFile` encrypt a `*wallet.Wallet` — one derived key pair. This CLI needs to encrypt something one level higher up the tree: the **mnemonic** itself, since that one phrase is what regenerates every derived address, forever, on demand (Chapter 38's whole point). The fix is a close cousin of Chapter 40's functions, reusing the exact same `deriveKeyFromPassword` (scrypt) and AES-256-GCM machinery, just wrapped around a plain string instead of a JSON-encoded key pair:

```go
package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"os"
)

// EncryptMnemonicFile encrypts a BIP-39 mnemonic phrase using a
// password-derived key and writes salt || nonce || ciphertext to path --
// the same on-disk layout and the same scrypt + AES-256-GCM machinery
// Chapter 40's EncryptWalletFile uses, applied one level higher: instead
// of protecting a single derived key pair, this protects the ONE secret
// an entire tree of addresses regenerates from.
func EncryptMnemonicFile(path, mnemonic, password string) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key, err := deriveKeyFromPassword(password, salt)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(mnemonic), nil)

	out := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	return os.WriteFile(path, out, 0600)
}

// DecryptMnemonicFile reverses EncryptMnemonicFile, returning the
// original mnemonic string if password is correct and the file has not
// been tampered with -- otherwise a generic error, for the exact same
// "don't hand an attacker a guessing oracle" reason Chapter 40 explained.
func DecryptMnemonicFile(path, password string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if len(data) < 16+12 {
		return "", errors.New("wallet: encrypted file is too short to be valid")
	}

	salt, nonce, ciphertext := data[:16], data[16:28], data[28:]

	key, err := deriveKeyFromPassword(password, salt)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("wallet: decryption failed (wrong password or corrupted file)")
	}

	return string(plaintext), nil
}
```

Notice how little of this is genuinely new: `deriveKeyFromPassword` is Chapter 40's exact scrypt wrapper, untouched. The AES-GCM seal/open logic, the `salt || nonce || ciphertext` layout, and the deliberately generic decryption error are all copied directly from `EncryptWalletFile`/`DecryptWalletFile`. The only real difference is *what* gets encrypted — a plain mnemonic string instead of a JSON-serialized key pair — which is exactly the point of building good, composable cryptographic primitives in earlier chapters: this chapter barely has to think about cryptography at all, because Chapter 40 already did that thinking once, correctly.

---

## 4. Reading Passwords Without Echoing Them

Chapter 36's CLI never actually prompted for a password (it had no encryption yet to protect with one). This CLI does, and typing a password into a terminal that echoes every character back to the screen — visible to anyone glancing at your monitor, or preserved in a terminal scrollback buffer — defeats much of the point of having a password at all. Go's `golang.org/x/term` package provides exactly the fix real command-line tools (`ssh`, `sudo`, `git`) use: reading a line of input from the terminal without echoing it back.

```go
package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

// promptPassword prints prompt, reads one line from the terminal WITHOUT
// echoing it to the screen, and returns it as a string. term.ReadPassword
// temporarily puts the terminal into "raw" mode for the read, restoring
// normal mode afterward automatically.
func promptPassword(prompt string) string {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // ReadPassword eats the Enter keystroke's newline; print our own
	if err != nil {
		log.Fatalf("reading password: %v", err)
	}
	return string(bytePassword)
}

// promptNewPassword asks for a password twice and requires both entries
// to match -- the same "confirm your password" pattern every sign-up
// form uses, applied here so a typo while choosing a NEW password can't
// silently lock a user out of their own just-created wallet file.
func promptNewPassword() string {
	for {
		p1 := promptPassword("Choose a password to encrypt this wallet: ")
		p2 := promptPassword("Confirm password: ")
		if p1 == p2 {
			return p1
		}
		fmt.Println("Passwords did not match -- try again.")
	}
}
```

`promptPassword` is the primitive: print whatever prompt the caller wants, then call `term.ReadPassword`, which reads a line directly from the terminal device in a mode where typed characters are captured but never displayed. `promptNewPassword` builds a slightly friendlier flow on top, specifically for the *first* time a password is chosen (`new` and `recover` both need this), where a silent typo would be far more costly than during an ordinary unlock (where a wrong password just fails cleanly and lets you try again).

---

## 5. Implementing new and recover

```go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/you/gochain/wallet"
)

func runNew(args []string) {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	keystorePath := fs.String("keystore", defaultKeystorePath, "path to write the encrypted wallet file")
	count := fs.Int("count", 5, "how many addresses to derive and print")
	fs.Parse(args)

	mnemonic, err := wallet.NewMnemonic()
	if err != nil {
		log.Fatalf("new: generating mnemonic: %v", err)
	}

	fmt.Println("=========================================================")
	fmt.Println(" YOUR NEW SEED PHRASE -- WRITE THIS DOWN. NEVER SHARE IT.")
	fmt.Println(" Anyone who has these words controls every coin every")
	fmt.Println(" address derived from this wallet will ever hold.")
	fmt.Println("=========================================================")
	fmt.Println(mnemonic)
	fmt.Println("=========================================================")

	password := promptNewPassword()

	if err := wallet.EncryptMnemonicFile(*keystorePath, mnemonic, password); err != nil {
		log.Fatalf("new: saving encrypted wallet: %v", err)
	}
	fmt.Printf("\nEncrypted wallet saved to %s\n\n", *keystorePath)

	printAddresses(mnemonic, *count)
}

func runRecover(args []string) {
	fs := flag.NewFlagSet("recover", flag.ExitOnError)
	keystorePath := fs.String("keystore", defaultKeystorePath, "path to write the recovered wallet file")
	count := fs.Int("count", 5, "how many addresses to derive and print")
	fs.Parse(args)

	fmt.Print("Enter your seed phrase: ")
	reader := bufio.NewReader(os.Stdin)
	mnemonic, _ := reader.ReadString('\n')
	mnemonic = strings.TrimSpace(mnemonic)

	if mnemonic == "" {
		log.Fatal("recover: no seed phrase entered")
	}

	password := promptNewPassword()

	if err := wallet.EncryptMnemonicFile(*keystorePath, mnemonic, password); err != nil {
		log.Fatalf("recover: saving encrypted wallet: %v", err)
	}
	fmt.Printf("\nWallet recovered and encrypted to %s\n\n", *keystorePath)

	printAddresses(mnemonic, *count)
}

// printAddresses derives and prints the first `count` addresses from a
// mnemonic -- shared by new, recover, and addresses, since all three
// ultimately need to show the same thing: "here is what this seed
// controls."
func printAddresses(mnemonic string, count int) {
	seed := wallet.SeedFromMnemonic(mnemonic, "")
	hd := wallet.NewHDWallet(seed)

	fmt.Println("Derived addresses:")
	for i := uint32(0); i < uint32(count); i++ {
		w, err := hd.DeriveWallet(i)
		if err != nil {
			log.Fatalf("deriving address %d: %v", i, err)
		}
		fmt.Printf("  [%d] %s\n", i, w.Address())
	}
}
```

`runNew` is a direct pipeline through this volume's own chapters, in order: `wallet.NewMnemonic()` (Chapter 39) generates the phrase, `promptNewPassword()` (Section 4) gets a password to protect it with, `wallet.EncryptMnemonicFile` (Section 3) persists it safely, and `printAddresses` proves the whole tree is immediately usable by deriving and printing real addresses from it. `runRecover` is deliberately almost identical, differing only in *where* the mnemonic comes from — typed in by the user via `bufio.Reader` instead of freshly generated — which is exactly the point: recovery should feel like an alternate entry point into the same pipeline, not a separate code path prone to drifting out of sync with `new` over time.

---

## 6. Implementing addresses, balance, and send

```go
package main

import (
	"encoding/hex"
	"flag"
	"log"

	"github.com/you/gochain/core"
	"github.com/you/gochain/wallet"
)

func runAddresses(args []string) {
	fs := flag.NewFlagSet("addresses", flag.ExitOnError)
	keystorePath := fs.String("keystore", defaultKeystorePath, "path to the encrypted wallet file")
	count := fs.Int("count", 5, "how many addresses to derive and print")
	fs.Parse(args)

	mnemonic := unlockMnemonic(*keystorePath)
	printAddresses(mnemonic, *count)
}

func runBalance(args []string) {
	fs := flag.NewFlagSet("balance", flag.ExitOnError)
	keystorePath := fs.String("keystore", defaultKeystorePath, "path to the encrypted wallet file")
	index := fs.Uint("index", 0, "which derived address index to check")
	fs.Parse(args)

	mnemonic := unlockMnemonic(*keystorePath)
	w := deriveAt(mnemonic, uint32(*index))

	bc := localChain()
	balance := bc.BalanceOf(w.Address())
	fmt.Printf("Address [%d] %s\nBalance: %d gochips\n", *index, w.Address(), balance)
}

func runSend(args []string) {
	fs := flag.NewFlagSet("send", flag.ExitOnError)
	keystorePath := fs.String("keystore", defaultKeystorePath, "path to the encrypted wallet file")
	fromIndex := fs.Uint("from-index", 0, "which derived address to send from")
	to := fs.String("to", "", "recipient address")
	amount := fs.Int64("amount", 0, "amount to send, in gochips")
	fs.Parse(args)

	if *to == "" || *amount <= 0 {
		log.Fatal("send: -to and a positive -amount are required")
	}

	mnemonic := unlockMnemonic(*keystorePath)
	w := deriveAt(mnemonic, uint32(*fromIndex))

	bc := localChain()
	utxoSet := core.NewUTXOSet(bc)

	tx, err := core.NewTransaction(w, *to, *amount, utxoSet)
	if err != nil {
		log.Fatalf("send: building transaction: %v", err)
	}

	prevTXs := findPrevTXs(bc, tx)
	tx.Sign(w.KeyPair.PrivateKey, prevTXs)

	mempool := core.NewMempool()
	if err := mempool.Add(tx); err != nil {
		log.Fatalf("send: transaction rejected by mempool: %v", err)
	}

	fmt.Printf("Transaction %x submitted to the mempool.\n", tx.ID)
}

// unlockMnemonic prompts for a password and decrypts the wallet file,
// exiting the program with a clear error on a wrong password rather
// than silently proceeding with garbage key material.
func unlockMnemonic(path string) string {
	password := promptPassword("Wallet password: ")
	mnemonic, err := wallet.DecryptMnemonicFile(path, password)
	if err != nil {
		log.Fatalf("could not unlock wallet: %v", err)
	}
	return mnemonic
}

// deriveAt is a small convenience wrapper: mnemonic -> seed -> HDWallet
// -> one specific derived *wallet.Wallet, the same three-step pipeline
// Chapter 40 built, used here every time a command needs one specific
// address's real key pair.
func deriveAt(mnemonic string, index uint32) *wallet.Wallet {
	seed := wallet.SeedFromMnemonic(mnemonic, "")
	hd := wallet.NewHDWallet(seed)
	w, err := hd.DeriveWallet(index)
	if err != nil {
		log.Fatalf("deriving address %d: %v", index, err)
	}
	return w
}

// localChain stands in for a real, shared, persisted chain. This mini
// project is about the WALLET half of GoChain specifically; Volumes 7-8
// are what eventually make "the chain" something genuinely shared
// across independent machines rather than a fresh, empty, in-process
// chain each invocation -- the exact same simplification Chapter 36's
// CLI made, for the same reasons.
func localChain() *core.Blockchain {
	return core.NewBlockchain()
}

func findPrevTXs(bc *core.Blockchain, tx *core.Transaction) map[string]core.Transaction {
	prevTXs := make(map[string]core.Transaction)
	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.TxID)
		if err != nil {
			log.Fatalf("could not find previous transaction %x: %v", in.TxID, err)
		}
		prevTXs[hex.EncodeToString(in.TxID)] = *prevTX
	}
	return prevTXs
}
```

`runAddresses` is the shortest command in the whole CLI, precisely because `unlockMnemonic` and `printAddresses` already do all the real work — a good sign that this chapter's design correctly factored out the shared logic. `runBalance` and `runSend` both follow the identical pattern Chapter 36 established: unlock (here, decrypt-and-derive instead of a flat file read), then hand off to the exact same `core.BalanceOf` and `core.NewTransaction`/`Transaction.Sign`/`mempool.Add` pipeline unchanged from Chapters 32-34 and Chapter 37's major project. `deriveAt` and `localChain` are the two seams where this chapter's HD-specific machinery meets the rest of GoChain — everything on the other side of those two functions has no idea, and does not need to care, whether the key pair it's holding came from a single flat file (Chapter 36) or from walking a whole derivation tree (this chapter).

---

## Mini Project: gochain-wallet

Here is the complete `main.go`, assembling every piece from Sections 3 through 6 into one runnable program (imports across the snippets above are combined at the top, as they would be in the real file):

```go
// cmd/gochain-wallet/main.go
package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/you/gochain/core"
	"github.com/you/gochain/wallet"
)

const defaultKeystorePath = "gochain-wallet.dat"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "new":
		runNew(os.Args[2:])
	case "recover":
		runRecover(os.Args[2:])
	case "addresses":
		runAddresses(os.Args[2:])
	case "balance":
		runBalance(os.Args[2:])
	case "send":
		runSend(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("usage: gochain-wallet <new|recover|addresses|balance|send> [flags]")
}

// -- runNew, runRecover, printAddresses: Section 5 --
// -- runAddresses, runBalance, runSend: Section 6 --
// -- unlockMnemonic, deriveAt, localChain, findPrevTXs: Section 6 --
// -- promptPassword, promptNewPassword: Section 4 --
```

A full sample terminal session — generating a wallet, "losing" access to the running program, and recovering the exact same wallet from its seed phrase alone — looks like this (the password is typed but never echoed, shown here as blank for realism):

```
$ go build -o gochain-wallet ./cmd/gochain-wallet

$ ./gochain-wallet new -count 3
=========================================================
 YOUR NEW SEED PHRASE -- WRITE THIS DOWN. NEVER SHARE IT.
 Anyone who has these words controls every coin every
 address derived from this wallet will ever hold.
=========================================================
witch collapse practice feed shame open maple north rescue fee coyote crumble
=========================================================
Choose a password to encrypt this wallet: 
Confirm password: 

Encrypted wallet saved to gochain-wallet.dat

Derived addresses:
  [0] 1Gz8kR9nQmWvTxYbLp3AeFj7HcRq2SdUo8
  [1] 1PxT4vLwZsKmNq9RyUeAg5BhCj8FoDm3Xz
  [2] 1Kj2mNp7AxLzWfHb6RtVcYq4EoGs9DiUp1

$ rm gochain-wallet.dat
$ echo "simulating a lost device -- the file above is now gone"

$ ./gochain-wallet recover -count 3
Enter your seed phrase: witch collapse practice feed shame open maple north rescue fee coyote crumble
Choose a password to encrypt this wallet: 
Confirm password: 

Wallet recovered and encrypted to gochain-wallet.dat

Derived addresses:
  [0] 1Gz8kR9nQmWvTxYbLp3AeFj7HcRq2SdUo8
  [1] 1PxT4vLwZsKmNq9RyUeAg5BhCj8FoDm3Xz
  [2] 1Kj2mNp7AxLzWfHb6RtVcYq4EoGs9DiUp1

$ ./gochain-wallet addresses -count 3
Wallet password: 
Derived addresses:
  [0] 1Gz8kR9nQmWvTxYbLp3AeFj7HcRq2SdUo8
  [1] 1PxT4vLwZsKmNq9RyUeAg5BhCj8FoDm3Xz
  [2] 1Kj2mNp7AxLzWfHb6RtVcYq4EoGs9DiUp1

$ ./gochain-wallet balance -index 0
Wallet password: 
Address [0] 1Gz8kR9nQmWvTxYbLp3AeFj7HcRq2SdUo8
Balance: 0 gochips
```

Notice the addresses printed after `recover` are byte-for-byte identical to the ones `new` printed originally — the entire promise of Chapter 38's hierarchical deterministic design, now demonstrated with a real, runnable tool rather than a diagram. Deleting the keystore file and typing the twelve words back in was enough to regenerate the exact same tree of addresses; the file itself was never the real backup, exactly as Chapter 40 argued.

---

## 8. What's Still Missing

- **`balance` and `send` still operate on a fresh, empty, in-memory chain per invocation** (via `localChain`), for the same reason Chapter 36's CLI did — there is no persistence (Volume 8) or real shared network (Volume 7) yet for this CLI to connect to. `balance -index 0` printing `0 gochips` immediately after generating a wallet is expected and correct: nobody has mined anything to that address on *this particular run's* throwaway chain.
- **No passphrase support at the CLI level.** `wallet.SeedFromMnemonic` accepts an optional passphrase (Chapter 38's "25th word"), but this CLI always passes an empty one. A production tool would add a `-passphrase` flag (prompted the same way as a password, via `promptPassword`) so users who want that extra layer can opt into it.
- **`NumAddresses` isn't tracked anywhere.** Every command that lists or uses addresses takes an explicit `-count` or `-index` flag; nothing here remembers "the user has already used addresses 0 through 4" between invocations. A polished version would persist a small amount of non-secret metadata (just the count, never the mnemonic) alongside the encrypted file.
- **No `mine` command**, matching Chapter 36's own scope note — this CLI can generate, recover, list, check, and send, but building blocks is left to the exercises, exactly as it was for Chapter 36.

---

## Summary

- This chapter rebuilds Chapter 36's `gochain-wallet` CLI on top of this volume's real HD wallet machinery: one seed phrase backing an entire tree of addresses, encrypted at rest instead of stored in plaintext.
- Five subcommands cover the whole lifecycle: `new` (generate), `recover` (restore from an existing phrase), `addresses` (list derived addresses), `balance` (check one), and `send` (spend from one) — mirroring Chapter 36's original three plus two new HD-specific ones.
- `EncryptMnemonicFile`/`DecryptMnemonicFile` are new this chapter, but reuse Chapter 40's `deriveKeyFromPassword` (scrypt) and AES-256-GCM machinery unchanged, applied to a plain mnemonic string instead of a serialized key pair.
- `golang.org/x/term.ReadPassword` reads a password from the terminal without echoing it to the screen — the same technique real CLI tools like `ssh` and `sudo` use, and a meaningful security improvement over Chapter 36's CLI, which never handled passwords at all.
- `promptNewPassword` requires a password to be typed twice and matched, specifically to avoid a silent typo locking a user out of a wallet they just created.
- A full sample session demonstrated the entire HD wallet promise concretely: generate a wallet, delete its keystore file entirely, recover it from the twelve-word phrase alone, and confirm every derived address regenerates byte-for-byte identical to the original.
- `balance` and `send` still operate against a fresh, local, in-memory chain per run — this mini project's scope is deliberately the wallet half of GoChain, not the shared-network half Volumes 7-8 build.

---

## Exercises

### Easy

1. Run through the full sample session in Section 7 on your own machine, using your own generated seed phrase. Confirm your own recovered addresses match your own originally-generated ones exactly.
2. Add a `-passphrase` flag to `new`, `recover`, `addresses`, `balance`, and `send`, prompted via `promptPassword` (labeled clearly as "optional extra passphrase, leave blank for none"), and pass it through to `wallet.SeedFromMnemonic` instead of the hardcoded empty string. Confirm that supplying a non-empty passphrase produces a completely different set of derived addresses than leaving it blank.
3. Modify `runAddresses` to also print each address's current index-derived public key hash (not just its human-readable address), using `wallet.HashPubKey` on the derived wallet's `KeyPair.PublicKey`, to make the connection between "address" and "public key hash" from Chapter 32 concrete.

### Medium

4. Implement the missing metadata tracking from Section 8: add a small, separate, unencrypted file (e.g. `gochain-wallet.meta.json`) alongside the encrypted keystore that records how many addresses have been shown to the user so far, and have `new`/`recover` initialize it and `addresses` update it (advancing it only if `-count` requests more addresses than previously recorded, never shrinking it).
5. Add a `change-password` subcommand that decrypts the wallet file with an old password and re-encrypts the same mnemonic with a new one (using `promptPassword` for the old password and `promptNewPassword` for the new one), without ever writing the plaintext mnemonic to disk at any intermediate step. Write a test proving the mnemonic recovered after a password change is identical to the mnemonic before it.
6. Currently, `runSend` builds a fresh, empty `core.NewMempool()` every single invocation, meaning a submitted transaction is visible to nobody, including a later `balance` check within the same process's lifetime. Rebuild `send` and `balance` to share a single, package-level, persisted (even just gob-encoded to a file, foreshadowing Chapter 20's approach) chain and mempool, so that sending coins in one invocation and checking a balance in a later, separate invocation of the CLI can actually observe each other's effects.

### Hard

7. This CLI derives addresses starting from index 0 every time, with no way to mark an address as "already used" versus "freshly generated for this specific payment." Design and implement a `next-address` subcommand that always returns the lowest-indexed address that has never been passed to `balance` or `send` before (tracked via the metadata file from Exercise 4), mimicking how real wallets avoid reusing the same receiving address repeatedly for privacy reasons (Chapter 38, Section 1).
8. Add a `-yes`/non-interactive mode to `send` that accepts the password via an environment variable (e.g. `GOCHAIN_WALLET_PASSWORD`) instead of an interactive terminal prompt, explicitly for use in scripts and automated testing — but require a loud, explicit warning printed to stderr whenever this mode is used, explaining the security tradeoff (a password sitting in an environment variable or shell history is far easier for another process or a shell history file to leak than one typed interactively and never echoed).
9. Research how real hardware wallets (Chapter 41's conceptual territory) would change this CLI's design if `gochain-wallet` were adapted to delegate signing to an external device instead of holding the mnemonic and deriving keys locally at all. Sketch (in comments or a short design document, no need to fully implement it) what `wallet.Signer` interface (previewed in Chapter 41) `deriveAt`'s return value would need to satisfy, and which of this chapter's functions (`unlockMnemonic`, `deriveAt`, or neither) would need to change versus which could stay exactly as written.
