# Chapter 69: Major Project 3 — Deploy a Token Contract

Volume 9 opened with a promise: an escrow contract that releases funds without an intermediary, because every node re-executes the same code and agrees on the result. Nine chapters later, every piece that promise needed is sitting in `gochain/vm`: an interpreter (Chapter 62), a scripting model unifying transactions and contracts (Chapter 63), gas metering (Chapter 64), a token contract's operations (Chapter 65), persistent per-contract storage (Chapter 66), a hard lesson in what happens when state updates happen in the wrong order (Chapter 67), and a test suite rigorous enough to trust (Chapter 68). This chapter's job is to prove all of it actually works together: deploy the token contract for real, mint a supply, move tokens between two wallets through a properly signed transaction, and read the resulting balances straight out of storage — the same fundamental flow behind every ERC-20 token that has ever moved on Ethereum, running here on a chain you built yourself.

## Table of Contents

1. [Recap: Every Piece This Project Assembles](#1-recap-every-piece-this-project-assembles)
2. [The Plan: Mint, Transfer, Query](#2-the-plan-mint-transfer-query)
3. [Signed Contract Invocations: Extending Trust to Contract Calls](#3-signed-contract-invocations-extending-trust-to-contract-calls)
4. [Setting Up the Node: Storage, Contract Store, Deploying the Token](#4-setting-up-the-node-storage-contract-store-deploying-the-token)
5. [Two Wallets, One Token](#5-two-wallets-one-token)
6. [Minting the Initial Supply](#6-minting-the-initial-supply)
7. [Building and Verifying a Signed Transfer](#7-building-and-verifying-a-signed-transfer)
8. [Querying Balances Directly From Storage](#8-querying-balances-directly-from-storage)
9. [Major Project: Token Deployment Demo](#9-major-project-token-deployment-demo)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Recap: Every Piece This Project Assembles

Every dependency this project needs already exists, built in an earlier chapter, and none of it needs to change:

```
wallet.New()  ------------------> two independent key pairs + addresses  (Volume 6)
crypto.Sign/Verify  -------------> proving a transfer was really authorized (Volume 2)
storage.OpenBoltStore  ----------> a real, crash-safe file on disk         (Chapter 55)
vm.NewContractStore  ------------> per-contract, isolated key-value state (Chapter 66)
vm.Token (Mint/Transfer/BalanceOf) -> the token contract's operations     (Chapter 68)
```

Nothing here is new machinery. This chapter's actual contribution is the wiring: a `ContractInvocation` type that extends Chapter 33's "sign exactly what you are authorizing" idea from plain coin transfers to contract method calls, so a transfer of tokens is just as provably authorized as a transfer of native gochips ever was.

One deliberate substitution worth naming up front: Chapter 65 designed `mint`, `transfer`, and `balanceOf` as raw VM bytecode — `Mint`, `Transfer`, and `BuildBalanceOfScript`, each assembling and running an `[]Instruction` program directly. This project deploys against `vm.Token` instead — Chapter 68's Go-level implementation of that exact same behavior, built directly on `ContractStore` "with the same clarity Chapter 67 tested `VulnerableBank` and `FixedBank`," as Chapter 68 Section 2 put it. Both express identical mint/transfer/balanceOf semantics against the same underlying storage; this chapter picks `vm.Token` because a deployment demo, like a test suite, benefits from that same clarity — the signing and authorization story below is exactly as real either way, since it lives in `ContractInvocation`, not in which of the two equivalent implementations happens to run underneath it.

---

## 2. The Plan: Mint, Transfer, Query

```
   Wallet W1 (issuer)              Wallet W2 (recipient)
        |                                  |
        |  1. mint(W1, 1,000,000)          |
        v                                  |
  [ Token contract storage ]               |
   W1 -> 1,000,000                         |
        |                                  |
        |  2. SIGNED transfer(W1, W2, 250) |
        |     signature verified BEFORE    |
        |     the transfer is applied      |
        v                                  v
  [ Token contract storage ]
   W1 ->   999,750
   W2 ->       250
        |
        |  3. balanceOf(W1), balanceOf(W2)
        v
  read directly from ContractStore, bypassing Token's own methods entirely
  -- proving the numbers really are on disk, not just in the Token
     value's in-memory view of them
```

Step 3 matters more than it looks: querying through `Token.BalanceOf` only proves the *Go type* thinks the balance is right. Reading the same slot back through `ContractStore.Get` directly — the same interface built in Chapter 66, with no `Token` type standing between the demo and the raw bytes on disk — proves the balance genuinely persisted where Chapter 66 said it would, isolated under this one contract's address, exactly as promised.

---

## 3. Signed Contract Invocations: Extending Trust to Contract Calls

Chapter 33 signed a trimmed copy of a transaction so a UTXO's owner, and only that owner, could authorize spending it. A contract call needs the identical guarantee: only the token holder whose balance is being debited should be able to authorize a transfer out of it. `ContractInvocation` applies Chapter 33's exact idea — sign precisely what is being authorized, nothing more — to a method call instead of a coin transfer:

```go
// vm/contract_invocation.go
package vm

import (
	"encoding/binary"
	"errors"

	"github.com/you/gochain/crypto"
)

// ErrInvalidSignature means a ContractInvocation's signature does not
// match its own SignedData() under its own claimed SenderPubKey — the
// contract call is rejected before any contract method ever runs.
var ErrInvalidSignature = errors.New("vm: contract invocation signature is invalid")

// ContractInvocation is a signed instruction to call one method of one
// deployed contract, with one argument list — the contract-call
// equivalent of a signed core.Transaction from Chapter 33.
type ContractInvocation struct {
	ContractAddress string
	Method          string
	Args            [][]byte
	SenderPubKey    []byte
	Signature       []byte
}

// SignedData returns exactly the bytes the signature must cover: the
// contract address, method name, and every argument, each individually
// length-prefixed so two different (method, args) combinations can never
// hash or sign identically by accident.
func (ci *ContractInvocation) SignedData() []byte {
	var buf []byte
	appendField := func(b []byte) {
		length := make([]byte, 4)
		binary.BigEndian.PutUint32(length, uint32(len(b)))
		buf = append(buf, length...)
		buf = append(buf, b...)
	}
	appendField([]byte(ci.ContractAddress))
	appendField([]byte(ci.Method))
	for _, arg := range ci.Args {
		appendField(arg)
	}
	return buf
}

// Verify reports whether Signature was really produced, by the private
// key matching SenderPubKey, over exactly SignedData() — reusing
// crypto.Verify from Chapter 13 unchanged, the same primitive
// OpCheckSig itself calls into.
func (ci *ContractInvocation) Verify() bool {
	return crypto.Verify(ci.SenderPubKey, ci.SignedData(), ci.Signature)
}
```

A node receiving a `ContractInvocation` follows one non-negotiable rule: call `Verify()` first, and refuse to run a single opcode of the contract if it returns `false`. This mirrors exactly how `OpCheckSig` (Chapter 62) and block validation (Chapter 19) both work — verification always happens *before* anything derived from the unverified data is trusted.

---

## 4. Setting Up the Node: Storage, Contract Store, Deploying the Token

```go
// cmd/tokendemo/main.go — Section 1 of the full program
package main

import (
	"fmt"
	"log"

	"github.com/you/gochain/storage"
	"github.com/you/gochain/vm"
	"github.com/you/gochain/wallet"
)

func main() {
	// A real, on-disk, crash-safe store — the same BoltStore Chapter 55
	// built, reused here for contract storage exactly as Chapter 66 wired
	// it up.
	store, err := storage.OpenBoltStore("tokendemo.db")
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer store.Close()

	contractStore := vm.NewContractStore(store)

	// Deploy the token contract. In a real node this would be a signed
	// deployment transaction of its own (left as Exercise 8); for this
	// demo, deploying is a direct, local call, exactly the way Chapter 66
	// Section 7's DeployContract is meant to be used before any network
	// layer is involved at all.
	token := vm.NewToken("gochain-token-v1", contractStore)
	fmt.Println("Token contract deployed at:", token.Address)
```

---

## 5. Two Wallets, One Token

```go
	// --- Section 2: two independent wallets, from Volume 6 ---
	issuer := wallet.New()
	recipient := wallet.New()

	fmt.Println("Issuer wallet:   ", issuer.Address())
	fmt.Println("Recipient wallet:", recipient.Address())
```

Nothing about these two lines is new — `wallet.New()` has meant exactly this since Chapter 36's CLI wallet. What is new is what these two addresses are about to be used for: not spending UTXOs, but authorizing calls into a contract's storage.

---

## 6. Minting the Initial Supply

```go
	// --- Section 3: mint the initial supply to the issuer ---
	const initialSupply = 1_000_000
	if err := token.Mint(issuer.Address(), initialSupply); err != nil {
		log.Fatalf("mint: %v", err)
	}
	fmt.Printf("Minted %d GCT to issuer (%s)\n", initialSupply, issuer.Address())
```

Minting is deliberately not signature-gated in this demo — a real deployment would restrict `mint` to whichever address deployed the contract (an access-control check, Chapter 67 Section 9's closing warning about missing permission checks), left here as Exercise 6 rather than folded into this project's main flow, to keep the signed-transfer path (this project's actual point) the center of attention.

---

## 7. Building and Verifying a Signed Transfer

```go
	// --- Section 4: a signed transfer, issuer -> recipient ---
	const transferAmount = 250

	invocation := &vm.ContractInvocation{
		ContractAddress: token.Address,
		Method:           "transfer",
		Args:             [][]byte{[]byte(issuer.Address()), []byte(recipient.Address())},
		SenderPubKey:     issuer.PublicKey(),
	}
	invocation.Signature = issuer.Sign(invocation.SignedData())

	// A node would receive exactly this struct over the network (Volume
	// 7) or through the API (Volume 10). Verification happens BEFORE a
	// single token moves — the same rule OpCheckSig and block validation
	// both already follow.
	if !invocation.Verify() {
		log.Fatal("contract invocation signature verification FAILED — refusing to run the transfer")
	}
	fmt.Println("Signed transfer verified. Applying: transfer", transferAmount, "GCT")

	if err := token.Transfer(issuer.Address(), recipient.Address(), transferAmount); err != nil {
		log.Fatalf("transfer: %v", err)
	}
```

A tampered invocation — a copied signature attached to a different amount, or a different recipient address — fails `Verify()` and never reaches `token.Transfer` at all, for exactly the same reason a corrupted transaction fails `crypto.Verify` in Chapter 13's worked example: the signature is bound to the *exact* bytes `SignedData()` produces, and changing even one byte of the contract address, method, or arguments changes what those bytes are.

---

## 8. Querying Balances Directly From Storage

```go
	// --- Section 5: query balances two ways ---
	fmt.Println("\nBalances via Token.BalanceOf:")
	fmt.Printf("  issuer:    %d GCT\n", token.BalanceOf(issuer.Address()))
	fmt.Printf("  recipient: %d GCT\n", token.BalanceOf(recipient.Address()))

	// Now read the SAME slots directly through ContractStore, bypassing
	// Token entirely — proving the balances really live on disk, isolated
	// under this contract's address, exactly as Chapter 66 designed.
	issuerRaw, _ := contractStore.Get(token.Address, []byte(issuer.Address()))
	recipientRaw, _ := contractStore.Get(token.Address, []byte(recipient.Address()))

	fmt.Println("\nSame balances, read directly from contract storage:")
	fmt.Printf("  issuer:    %d GCT (raw slot bytes: %x)\n", vm.BytesToUint64(issuerRaw), issuerRaw)
	fmt.Printf("  recipient: %d GCT (raw slot bytes: %x)\n", vm.BytesToUint64(recipientRaw), recipientRaw)
}
```

```go
// vm/encoding.go (this chapter's addition)

// BytesToUint64 is an exported wrapper around the package-private
// bytesToUint64 helper Chapter 62 defined for internal opcode use — added
// here purely so code outside gochain/vm, like cmd/tokendemo (which,
// being package main, cannot see unexported names in another package),
// can decode a raw storage slot's bytes for display. The encoding itself
// is unchanged: the same big-endian uint64 format fixed since Chapter 62.
func BytesToUint64(b []byte) uint64 {
	return bytesToUint64(b)
}
```

---

## 9. Major Project: Token Deployment Demo

Putting Sections 4 through 8 together into one file gives the complete, runnable project:

```go
// cmd/tokendemo/main.go (complete)
package main

import (
	"fmt"
	"log"

	"github.com/you/gochain/storage"
	"github.com/you/gochain/vm"
	"github.com/you/gochain/wallet"
)

func main() {
	store, err := storage.OpenBoltStore("tokendemo.db")
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer store.Close()

	contractStore := vm.NewContractStore(store)
	token := vm.NewToken("gochain-token-v1", contractStore)
	fmt.Println("Token contract deployed at:", token.Address)

	issuer := wallet.New()
	recipient := wallet.New()
	fmt.Println("Issuer wallet:   ", issuer.Address())
	fmt.Println("Recipient wallet:", recipient.Address())

	const initialSupply = 1_000_000
	if err := token.Mint(issuer.Address(), initialSupply); err != nil {
		log.Fatalf("mint: %v", err)
	}
	fmt.Printf("Minted %d GCT to issuer (%s)\n", initialSupply, issuer.Address())

	const transferAmount = 250
	invocation := &vm.ContractInvocation{
		ContractAddress: token.Address,
		Method:          "transfer",
		Args:            [][]byte{[]byte(issuer.Address()), []byte(recipient.Address())},
		SenderPubKey:    issuer.PublicKey(),
	}
	invocation.Signature = issuer.Sign(invocation.SignedData())

	if !invocation.Verify() {
		log.Fatal("contract invocation signature verification FAILED — refusing to run the transfer")
	}
	fmt.Println("Signed transfer verified. Applying: transfer", transferAmount, "GCT")

	if err := token.Transfer(issuer.Address(), recipient.Address(), transferAmount); err != nil {
		log.Fatalf("transfer: %v", err)
	}

	fmt.Println("\nBalances via Token.BalanceOf:")
	fmt.Printf("  issuer:    %d GCT\n", token.BalanceOf(issuer.Address()))
	fmt.Printf("  recipient: %d GCT\n", token.BalanceOf(recipient.Address()))

	issuerRaw, _ := contractStore.Get(token.Address, []byte(issuer.Address()))
	recipientRaw, _ := contractStore.Get(token.Address, []byte(recipient.Address()))
	fmt.Println("\nSame balances, read directly from contract storage:")
	fmt.Printf("  issuer:    %d GCT (raw slot bytes: %x)\n", vm.BytesToUint64(issuerRaw), issuerRaw)
	fmt.Printf("  recipient: %d GCT (raw slot bytes: %x)\n", vm.BytesToUint64(recipientRaw), recipientRaw)
}
```

```
$ go run ./cmd/tokendemo
Token contract deployed at: gochain-token-v1
Issuer wallet:    1A2b3C...W1addr
Recipient wallet: 9Z8y7X...W2addr
Minted 1000000 GCT to issuer (1A2b3C...W1addr)
Signed transfer verified. Applying: transfer 250 GCT

Balances via Token.BalanceOf:
  issuer:    999750 GCT
  recipient: 250 GCT

Same balances, read directly from contract storage:
  issuer:    999750 GCT (raw slot bytes: 00000000000f423e)
  recipient: 250 GCT (raw slot bytes: 00000000000000fa)
```

Run it a second time without deleting `tokendemo.db`, and minting adds another 1,000,000 to the issuer's already-existing balance — the same lesson Chapter 55's `storagedemo` closed with: the store alone does not stop you from calling an operation twice; it is the calling code's job (here, and in a real node, some higher-level "has this deployment already happened" check) to decide when an operation should run at all. This project's actual contribution is proving that once a signed, verified call does run, the resulting state is real, on disk, isolated per contract, and readable two independent ways that agree with each other exactly.

---

## Summary

- This project wires together every piece Volume 9 built: the VM and its opcodes (Ch 62), gas metering (Ch 64), the token contract's operations (Ch 65), persistent per-contract storage (Ch 66), the reentrancy fix (Ch 67), and the test suite that holds it all to account (Ch 68).
- `ContractInvocation` extends Chapter 33's "sign exactly what you are authorizing" rule from plain coin transfers to contract method calls: a length-prefixed encoding of the contract address, method, and arguments, signed with the caller's private key.
- `Verify()` must run — and pass — before a single opcode of a contract call executes, mirroring `OpCheckSig` and block validation's shared rule that verification always precedes trust.
- The demo mints an initial supply to one wallet, then moves tokens to a second wallet only after checking a real ECDSA signature over the exact transfer being requested.
- Balances are queried two independent ways — through `Token.BalanceOf` and directly through `ContractStore.Get` — and agree, proving the state genuinely persisted on disk rather than only existing in one in-memory view of it.
- Minting in this demo is deliberately left without an access-control check, flagged explicitly as a simplification (and an exercise) rather than a silent gap.
- Running the demo twice without resetting the database compounds the mint, the same lesson Chapter 55 closed its own hands-on demo with: storage alone does not enforce business rules — the calling code still has to.

---

## Exercises

### Easy

1. What specific bytes does `ContractInvocation.SignedData()` include, in order? If two invocations had identical `Method` and `Args` but different `ContractAddress` values, would they ever produce the same signed data? Why does that matter?
2. Trace through what happens if `invocation.Verify()` is accidentally called *after* `token.Transfer` instead of before it, in a version of this demo where the invocation happens to be invalid. What real-world harm does checking order avoid here, connecting back to Chapter 67 Section 3's checks-effects-interactions rule?
3. The demo prints balances two ways in Section 8. If `Token.BalanceOf(issuer.Address())` and the raw `contractStore.Get(...)` read ever disagreed with each other, what would that tell you about a bug in `Token`'s implementation specifically (as opposed to a bug in `ContractStore`)?

### Medium

4. Modify the demo so a *third* wallet attempts to submit a `ContractInvocation` claiming to transfer from the issuer's address, but signs it with its own private key instead of the issuer's. Confirm `Verify()` correctly rejects it, and print a clear message distinguishing this failure from an insufficient-balance failure.
5. Add a `MaxSupply` field to `Token` and a check in `Mint` that refuses to mint past it, returning a new `ErrMaxSupplyExceeded`. Update the demo to attempt minting an amount that would exceed a deliberately small `MaxSupply`, and show the mint being correctly refused.
6. This chapter's `Mint` call has no access-control check at all — anyone who can call it can mint unlimited tokens. Add an `Owner` field to `Token`, set at construction, and modify `Mint` to require a `ContractInvocation`-style signed call from the owner specifically, rejecting any other caller's mint attempt. Write a test proving a non-owner's mint is rejected.
7. Extend the demo to perform a *second* transfer, recipient-to-issuer, sending half of what the recipient just received back, and print a running balance history showing all three states (initial mint, after transfer 1, after transfer 2).

### Hard

8. Design (in prose, plus a sketch of the Go types involved) how contract *deployment* itself — not just calls to an already-deployed contract — could be made into a signed operation, so a node only accepts a new contract from whoever's signature accompanies its bytecode. How would this interact with `DeployContract`'s address derivation from Chapter 66 (which currently depends only on the code, not on who deployed it)?
9. This demo runs entirely inside one Go process, with no networking involved at all. Sketch how this same flow would look across Volume 7's P2P network instead: which node verifies the `ContractInvocation`'s signature, when does the resulting state change propagate to other nodes, and what happens if two conflicting, both-individually-valid transfers from the same balance arrive at two different nodes nearly simultaneously? (You do not need to implement this — a clear, reasoned design, referencing specific mechanisms from Volume 7, is the goal.)
