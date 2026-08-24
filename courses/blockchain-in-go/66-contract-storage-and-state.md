# Chapter 66: Contract Storage and State

Every opcode built in Chapter 62 does its work on the stack and then forgets it ever happened. That is fine for `OpAdd` — nobody needs `5 + 3` to still equal `8` tomorrow — but it is fatal for a token contract. The whole point of the `mint`, `transfer`, and `balanceOf` operations Chapter 65 designed is that a balance set today must still be there the next time anyone calls the contract, possibly seconds later, possibly after the node has restarted entirely. `OpSLoad` and `OpSStore` have had the right stack effects since Chapter 62 — `( key -- value )` and `( value key -- )` — but their bodies were placeholders: reads always came back empty, writes vanished into nothing. This chapter replaces those placeholders with a real, persistent, per-contract storage layer, and gives the VM a `Contract` type that can actually be deployed and called.

## Table of Contents

1. [Recap: Two Stub Opcodes Waiting for a Real Backend](#1-recap-two-stub-opcodes-waiting-for-a-real-backend)
2. [Why a Contract Needs Its Own Storage, Not the UTXO Set](#2-why-a-contract-needs-its-own-storage-not-the-utxo-set)
3. [The Isolation Problem, and How Keys Solve It](#3-the-isolation-problem-and-how-keys-solve-it)
4. [Extending `storage.Store` With Generic Put/Get](#4-extending-storagestore-with-generic-putget)
5. [Implementing `ContractStore`](#5-implementing-contractstore)
6. [Wiring `OpSLoad`/`OpSStore` Into the VM](#6-wiring-opsloadopsstore-into-the-vm)
7. [`Contract`, `DeployContract`, and `Call`](#7-contract-deploycontract-and-call)
8. [Hands-On: Two Contracts, Provably Isolated Storage](#8-hands-on-two-contracts-provably-isolated-storage)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Recap: Two Stub Opcodes Waiting for a Real Backend

Chapter 62 left a clear note where this chapter's work begins:

```go
func (vm *VM) execSLoad() error {
	// key is popped and, for now, ignored — there is nowhere yet to look
	// it up. Chapter 66 replaces this body with a real lookup against a
	// contract's storage, keeping this exact stack effect: ( key -- value ).
	if _, err := vm.stack.Pop(); err != nil {
		return err
	}
	vm.stack.Push([]byte{})
	return nil
}
```

Nothing about the *shape* of these opcodes needs to change: a program that calls `OpSLoad` still expects to pop a key and get a value back; a program that calls `OpSStore` still expects to pop a value and a key and get nothing back. What changes is entirely inside the method body, and entirely invisible to any VM program written against the stack effect alone — which is exactly the point of having fixed that stack effect first, back in Chapter 61, before a single byte of Go existed.

---

## 2. Why a Contract Needs Its Own Storage, Not the UTXO Set

GoChain already has a place where value lives: the UTXO set from Volume 5 and Volume 8, tracking which unspent outputs belong to which address. It is tempting to ask why a contract cannot just reuse that. The answer is that a UTXO is a specific *shape* of data — an amount and an owner — and a contract's storage needs to hold whatever shape the contract's own logic requires: a token contract's balance table, yes, but just as easily a voting contract's tally, or an escrow contract's "has the buyer confirmed?" flag. The UTXO set answers one question ("who owns this coin"); contract storage needs to answer whatever question the contract was written to answer.

Think of the UTXO set as a bank's central ledger of cash in circulation, and contract storage as a filing cabinet sitting inside one particular business — a token issuer, a voting booth, an escrow office. Both live in the same building (the same GoChain node, backed by the same `storage.Store`), but a customer's cash balance and a business's internal paperwork are not the same kind of record, and mixing them into one table would make both harder to reason about.

```
                    ONE GoChain NODE, ONE BoltDB FILE
   +------------------------------------------------------------------+
   |  blocks bucket    utxo bucket        contracts bucket (this      |
   |  (Ch 55)          (Ch 55/56)          chapter, new)               |
   |                                                                    |
   |  "who owns       "how much does      "whatever each contract's    |
   |   this block"      this address       own code decided to        |
   |                     hold"              remember"                 |
   +------------------------------------------------------------------+
```

---

## 3. The Isolation Problem, and How Keys Solve It

Say GoChain ends up running two contracts: the token contract from Chapter 65, and some other contract a different developer deploys later — call it `Address B`. Both contracts might, entirely coincidentally, decide to use the byte string `"alice"` as a storage slot: the token contract to record Alice's token balance, contract B to record something unrelated, like Alice's vote in a poll. If both contracts' storage lived in one shared table keyed only by slot, `"alice"`, these would collide — contract B's write would silently overwrite the token contract's balance for Alice, or vice versa, depending on write order. Neither contract's author did anything wrong; they just happened to pick the same word.

The fix is the same one Chapter 14 used for addresses and Chapter 55 used for buckets: make the *real* key wider than what any single piece of code sees. Every storage slot a contract asks for gets silently combined with that contract's own address before it ever reaches disk, so `"alice"` under the token contract's address and `"alice"` under contract B's address land at two completely different physical locations — even though both contracts' code only ever mentions the bare slot `"alice"`.

```
Contract code sees:              Physical key that actually gets stored:

TOKEN CONTRACT (addr 7f3a...)
  SStore("alice", 100)     -->   [len=20][7f3a...............][alice]
  SLoad("alice")           -->   [len=20][7f3a...............][alice]

CONTRACT B      (addr 91c2...)
  SStore("alice", "yes")   -->   [len=20][91c2...............][alice]
  SLoad("alice")           -->   [len=20][91c2...............][alice]

Same slot name, "alice" — but the length-prefixed contract address in
front of it guarantees the two physical keys can never collide, no
matter what bytes either contract's address or slot happens to contain.
```

A length prefix in front of the address (rather than, say, a `:` separator) is deliberate: a `:` separator would break if a contract address or a slot name ever happened to contain a `:` byte itself, silently merging two logically different keys. A 4-byte length prefix removes that ambiguity entirely — there is exactly one way to split the combined key back into "address part" and "slot part," no matter what bytes either part contains.

---

## 4. Extending `storage.Store` With Generic Put/Get

Chapter 55 designed `storage.Store` around exactly what `core` and `consensus` needed at the time: blocks and UTXOs, each with their own narrow, typed methods. Contract storage needs something more general — arbitrary keys, arbitrary values, no fixed schema, because a contract's storage slots are whatever that contract's own code decides they are. Rather than inventing a seventh narrowly-typed method, this chapter adds the two most general operations a key-value store can offer:

```go
// storage/store.go (Chapter 66 addition — the five Chapter 55 methods are
// unchanged and simply not repeated here)
package storage

import "errors"

// ErrNotFound is returned by Get when key has never been written, so
// callers (like ContractStore, Section 5) can distinguish "this slot is
// legitimately empty" from "the database itself failed."
var ErrNotFound = errors.New("storage: not found")

type Store interface {
	PutBlock(hash []byte, block *core.Block) error
	GetBlock(hash []byte) (*core.Block, error)
	PutUTXO(key string, output *core.TxOutput) error
	GetUTXO(key string) (*core.TxOutput, error)
	DeleteUTXO(key string) error
	Iterator() Iterator

	// Put and Get are the fully generic pair contract storage depends on:
	// no schema, no decoding, just bytes in and bytes back out.
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
}
```

`BoltStore` gets a fourth bucket to back these two new methods, keeping the same "one drawer, one job" discipline Chapter 55 established for `blocks`, `utxo`, and `meta`:

```go
// storage/bolt_store.go (Chapter 66 addition)
var contractsBucket = []byte("contracts")

// (added to the loop of buckets OpenBoltStore creates via
// CreateBucketIfNotExists, alongside blocksBucket, utxoBucket, metaBucket)

func (s *BoltStore) Put(key []byte, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(contractsBucket).Put(key, value)
	})
}

func (s *BoltStore) Get(key []byte) ([]byte, error) {
	var value []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(contractsBucket).Get(key)
		if data == nil {
			return ErrNotFound
		}
		// Copy out, exactly as Chapter 55 Section 8 explained: data points
		// into a memory-mapped page that stops being valid the instant
		// this View transaction ends.
		value = append([]byte(nil), data...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}
```

Nothing about `blocks`, `utxo`, or `meta` changes. A fourth, independent bucket means a bug in contract storage code can no more corrupt the block log than a bug in the UTXO index can — the same argument Chapter 55 made for keeping those two separate applies again here, for a third time, without needing to re-litigate it.

---

## 5. Implementing `ContractStore`

With `storage.Store` now generic enough, `ContractStore` is a thin, deliberate layer on top of it whose entire job is Section 3's key-widening trick:

```go
// vm/contract_store.go
package vm

import (
	"encoding/binary"
	"errors"

	"github.com/you/gochain/storage"
)

// ContractStore persists per-contract key-value state, backed by
// gochain/storage.Store, keyed by (contractAddress + slotKey). No two
// contracts' storage can ever collide, no matter what slot names either
// one happens to choose.
type ContractStore struct {
	store storage.Store
}

// NewContractStore wraps an already-open storage.Store (typically a
// *storage.BoltStore, but any implementation works, per Chapter 55's
// interface-first design) for contract state specifically.
func NewContractStore(store storage.Store) *ContractStore {
	return &ContractStore{store: store}
}

// storageKey combines a contract's address and one of its storage slots
// into the single physical key that actually reaches the database. A
// 4-byte big-endian length prefix in front of the address means the split
// between "address part" and "slot part" is always unambiguous, even if
// either part happens to contain bytes that look like a separator.
func storageKey(contractAddr string, slot []byte) []byte {
	addr := []byte(contractAddr)
	key := make([]byte, 4+len(addr)+len(slot))
	binary.BigEndian.PutUint32(key[:4], uint32(len(addr)))
	copy(key[4:4+len(addr)], addr)
	copy(key[4+len(addr):], slot)
	return key
}

// Get reads one storage slot belonging to contractAddr. A slot that was
// never written comes back as an empty, non-nil value and a nil error —
// matching the exact behavior Chapter 62's stub gave OpSLoad, so no VM
// program needs to change how it handles "nothing here yet."
func (cs *ContractStore) Get(contractAddr string, slot []byte) ([]byte, error) {
	value, err := cs.store.Get(storageKey(contractAddr, slot))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return []byte{}, nil
		}
		return nil, err
	}
	return value, nil
}

// Set writes one storage slot belonging to contractAddr, overwriting
// whatever (if anything) was there before.
func (cs *ContractStore) Set(contractAddr string, slot []byte, value []byte) error {
	return cs.store.Put(storageKey(contractAddr, slot), value)
}
```

Every method here does exactly one thing: turn a `(contractAddr, slot)` pair into one physical key, then delegate straight to the generic `Store` built in Section 4. `ContractStore` itself holds no state of its own beyond the `Store` it wraps — restart the node, reopen the same BoltDB file, and every contract's storage is exactly where it left off, because it was never anywhere but on disk in the first place.

---

## 6. Wiring `OpSLoad`/`OpSStore` Into the VM

The VM needs to know two things to make a real storage call: which `ContractStore` to talk to, and which contract's address to store under during this particular run. Both are new, unexported fields on `VM` — adding them changes nothing about `NewVM`'s signature from Chapter 62, since a freshly constructed VM simply starts with no contract context at all, exactly as it starts with an empty stack:

```go
// vm/vm.go (Chapter 66 additions to the VM struct from Chapter 62)
type VM struct {
	stack    *Stack
	program  []Instruction
	pc       int
	gasUsed  uint64
	gasLimit uint64

	contractStore *ContractStore // which storage backend to read/write
	contractAddr  string         // whose storage slots this run reads/writes
}

// ErrNoContractStore means a program tried to run OpSLoad or OpSStore on a
// VM that was never given a ContractStore — a plain arithmetic program,
// with no contract behind it, should never hit this, but a contract
// program run without AttachStorage should fail loudly rather than
// silently behaving like Chapter 62's empty stub again.
var ErrNoContractStore = errors.New("vm: no contract storage attached to this vm")

// AttachStorage gives vm access to persistent, per-contract storage.
// Call it once, right after NewVM, before running any program that uses
// OpSLoad or OpSStore. A VM used only for arithmetic (Chapter 62's tests,
// for instance) never needs to call this at all.
func (vm *VM) AttachStorage(cs *ContractStore) {
	vm.contractStore = cs
}
```

And the two opcode bodies themselves, replacing Chapter 62's stubs line for line — same stack effects, same pop order, real behavior underneath:

```go
// vm/opcodes.go (Chapter 66 replaces these two method bodies)

func (vm *VM) execSLoad() error {
	key, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	if vm.contractStore == nil {
		return ErrNoContractStore
	}
	value, err := vm.contractStore.Get(vm.contractAddr, key)
	if err != nil {
		return err
	}
	vm.stack.Push(value)
	return nil
}

func (vm *VM) execSStore() error {
	key, err := vm.stack.Pop() // key is nearest the top, same as before
	if err != nil {
		return err
	}
	value, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	if vm.contractStore == nil {
		return ErrNoContractStore
	}
	return vm.contractStore.Set(vm.contractAddr, key, value)
}
```

Every test Chapter 62 wrote against these two opcodes' stack effects (`TestOpSLoad_EmptyByDefault`, `TestOpSStore_DoesNotError`) still passes unchanged, once the VM under test calls `AttachStorage` first — the observable contract (pop this many values, in this order, push this many back) never moved. Only the question "where does the value actually come from, and where does it actually go" got a real answer.

---

## 7. `Contract`, `DeployContract`, and `Call`

A contract is a named, deployed piece of code plus the address its storage lives under. `DeployContract` derives that address deterministically from the code itself — hashing a contract's bytecode the same way Chapter 14 hashed a public key into an address, so two nodes that both deploy byte-for-byte identical code independently arrive at the same address without needing to agree on anything else first:

```go
// vm/contract.go
package vm

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"

	"github.com/you/gochain/crypto"
)

// Contract is one deployed, callable piece of GoChain VM bytecode. Address
// is derived from Code, so identical code always deploys to the identical
// address, and Code is exactly what Execute will run on every Call.
type Contract struct {
	Address string
	Code    []Instruction
}

// DeployContract computes a contract's address from its own bytecode and
// returns a ready-to-call Contract. Deploying does not run any code —
// exactly like NewVM back in Chapter 62, constructing does not execute.
func DeployContract(code []Instruction) *Contract {
	return &Contract{
		Address: deriveContractAddress(code),
		Code:    code,
	}
}

func deriveContractAddress(code []Instruction) string {
	var buf bytes.Buffer
	// Encoding errors here would mean a fundamentally malformed
	// Instruction slice — not something a deploy step can recover from,
	// so a panic here is deliberate rather than an omission.
	if err := gob.NewEncoder(&buf).Encode(code); err != nil {
		panic(fmt.Sprintf("vm: cannot encode contract code: %v", err))
	}
	hash := crypto.Hash(buf.Bytes()) // built in Chapter 09
	return hex.EncodeToString(hash)[:40]
}

// Call runs one method of the contract against vm, which must already
// have a ContractStore attached (Section 6) — Call binds vm's contract
// context to c.Address for the duration of this run, so every OpSLoad and
// OpSStore the program executes reaches exactly this contract's slots and
// no one else's.
func (c *Contract) Call(vm *VM, method string, args [][]byte) ([]byte, error) {
	vm.contractAddr = c.Address

	// Arguments go on first (deepest), the method name goes on last
	// (nearest the top) — Chapter 65's dispatcher pops the method name
	// first to decide which branch of the contract's code to jump to,
	// then pops arguments in the order that branch expects them.
	for _, arg := range args {
		vm.stack.Push(arg)
	}
	vm.stack.Push([]byte(method))

	vm.program = c.Code
	vm.pc = 0

	if err := vm.Execute(); err != nil {
		return nil, fmt.Errorf("contract %s: call %q: %w", c.Address, method, err)
	}

	result, err := vm.stack.Pop()
	if err != nil {
		return nil, fmt.Errorf("contract %s: call %q: no return value: %w", c.Address, method, err)
	}
	return result, nil
}
```

`Call` deliberately reuses one `*VM` across many calls rather than constructing a fresh one every time: `vm.pc = 0` rewinds the program counter, `vm.program = c.Code` loads this contract's bytecode, and the stack is left exactly as `Execute` left it — empty, if the contract's own code balances every push with a pop, which a well-behaved contract always does. Gas accounting (Chapter 64) and storage (this chapter) both live on the same long-lived `VM`, so a single node process can call the same contract thousands of times without re-doing any setup beyond this.

---

## 8. Hands-On: Two Contracts, Provably Isolated Storage

Section 3 argued, in prose, that two contracts choosing the same slot name cannot collide. Here is that argument as a test — two completely independent programs, both writing to the slot `"alice"`, both read back correctly:

```go
// vm/contract_store_test.go
package vm

import "testing"

func TestContractStore_IsolatesByAddress(t *testing.T) {
	backing := newFakeStore() // a tiny in-memory storage.Store, for tests only
	cs := NewContractStore(backing)

	err1 := cs.Set("contractA-address", []byte("alice"), uint64ToBytes(100))
	err2 := cs.Set("contractB-address", []byte("alice"), []byte("yes"))
	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected errors: %v, %v", err1, err2)
	}

	gotA, err := cs.Get("contractA-address", []byte("alice"))
	if err != nil {
		t.Fatalf("get from contract A: %v", err)
	}
	if bytesToUint64(gotA) != 100 {
		t.Fatalf("contract A's \"alice\" slot: expected 100, got %d", bytesToUint64(gotA))
	}

	gotB, err := cs.Get("contractB-address", []byte("alice"))
	if err != nil {
		t.Fatalf("get from contract B: %v", err)
	}
	if string(gotB) != "yes" {
		t.Fatalf("contract B's \"alice\" slot: expected \"yes\", got %q", gotB)
	}
}

func TestContractStore_UnwrittenSlotIsEmpty(t *testing.T) {
	cs := NewContractStore(newFakeStore())
	value, err := cs.Get("some-contract", []byte("never-written"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(value) != 0 {
		t.Fatalf("expected empty value for an unwritten slot, got %v", value)
	}
}
```

And, end to end through the VM itself: deploy a contract whose whole job is "store whatever the caller says under the slot `balance`, then read it back," call it twice with two different `ContractStore`-backed values, and confirm gas and storage both behave exactly as Sections 6 and 7 promised:

```go
func TestContract_SStoreThenSLoad(t *testing.T) {
	code := []Instruction{
		{Op: OpPop},                                // discard the method name; this toy contract has only one operation
		{Op: OpPush, Arg: []byte("balance")},        // key
		{Op: OpSStore},                              // ( value key -- ), value was pushed by Call as an argument
		{Op: OpPush, Arg: []byte("balance")},
		{Op: OpSLoad},                               // ( key -- value )
		{Op: OpHalt},
	}
	contract := DeployContract(code)

	cs := NewContractStore(newFakeStore())
	v := NewVM(nil, 1000)
	v.AttachStorage(cs)

	result, err := contract.Call(v, "set", [][]byte{uint64ToBytes(42)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bytesToUint64(result) != 42 {
		t.Fatalf("expected 42 read back, got %d", bytesToUint64(result))
	}
}
```

`newFakeStore()` here is a minimal in-memory `storage.Store` used only in tests, satisfying the exact same interface Section 4 extended — Chapter 55's Section 2 argument about swappable implementations paying off directly: nothing about `ContractStore` or `Contract.Call` needed to know or care that this test never touched a real BoltDB file.

---

## Summary

- `OpSLoad` and `OpSStore` kept the exact stack effects Chapter 61 specified and Chapter 62 stubbed out; this chapter replaced only their method bodies with a real, persistent backend.
- `storage.Store` (Chapter 55) grew two fully generic methods, `Put`/`Get`, backed by a new `contracts` bucket in `BoltStore` — contract storage needs no fixed schema, unlike blocks or UTXOs.
- `ContractStore` combines a contract's address and a storage slot into one physical key, using a length-prefixed encoding that makes the address/slot split unambiguous no matter what bytes either one contains.
- That combined-key scheme is what makes two different contracts' storage provably non-colliding, even if both happen to choose the exact same slot name.
- `VM` gained two new fields — `contractStore` and `contractAddr` — set via `AttachStorage` and by `Contract.Call` respectively; `NewVM`'s signature from Chapter 62 did not need to change at all.
- `Contract` and `DeployContract` give GoChain a deployable unit: an address derived deterministically from bytecode, plus the bytecode itself.
- `Contract.Call` pushes arguments and a method name onto the stack, points the VM at the contract's code, and runs it — reusing one long-lived `*VM` across many calls rather than reconstructing one each time.

---

## Exercises

### Easy

1. Why does `storageKey` use a 4-byte length prefix instead of a separator byte like `:` between the contract address and the slot? Construct a concrete pair of address/slot values that would collide under the `:`-separator scheme but not under the length-prefixed one.
2. `ContractStore.Get` treats `storage.ErrNotFound` as "return an empty value, nil error" rather than propagating the error. Why does that match Chapter 62's original `execSLoad` stub behavior, and why would silently swallowing every kind of error (not just `ErrNotFound`) here be a mistake?
3. `AttachStorage` is called once per `*VM`, not once per `Contract.Call`. What would go wrong (be concrete) if `Contract.Call` created a brand-new `*VM` with `NewVM` on every single call instead of reusing one?

### Medium

4. Add a `Delete(contractAddr string, slot []byte) error` method to `ContractStore`, backed by a new `Delete(key []byte) error` method on `storage.Store` (and `BoltStore`). Write a test proving a deleted slot reads back as empty again, matching the "unwritten" case.
5. `deriveContractAddress` panics if gob encoding fails. Under what realistic circumstances (if any) could encoding a `[]Instruction` actually fail? Is a panic the right response here, or should `DeployContract` return an `error` instead? Argue for one and rewrite the signature if you choose the latter.
6. Write `TestContract_TwoCallsShareState` that calls a "counter" contract's `increment` method three times in a row on the same `*VM` and `ContractStore`, asserting the stored count is 3 after the third call — proving storage genuinely persists across separate `Call` invocations, not just within one.

### Hard

7. `Contract.Call` currently has no limit on how many times it can be called, or any way to run two calls to the *same* contract concurrently on two different `*VM` instances sharing one `ContractStore`. Write a test using two goroutines that both call `increment` on the same counter contract 100 times each, and determine (with evidence, not just a guess) whether the final stored count is reliably 200. If it is not, explain exactly where the race occurs and propose a fix.
8. Extend `storageKey` so it can also express "all slots belonging to contract X" as a queryable prefix (useful for a future "dump this contract's entire storage" debugging tool), and implement `func (cs *ContractStore) Keys(contractAddr string) ([][]byte, error)` on top of it, backed by a new `PrefixIterator` you add to `storage.Store`. Explain why the length-prefix scheme from Section 3 makes this safe to implement without accidentally matching another contract's slots.
9. This chapter derives a contract's address purely from its bytecode, meaning identical code always produces an identical address, even if deployed by two different, unrelated users. Is this a problem in practice? Research how real systems (Ethereum's `CREATE` vs. `CREATE2` opcodes) handle contract address derivation, and propose (in prose, plus a sketch of a modified `DeployContract` signature) a scheme where GoChain could let the same bytecode be deployed by different callers at different addresses.
