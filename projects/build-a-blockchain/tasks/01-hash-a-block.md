# Task 01: Hash a Block

## What you will build

The foundation everything else sits on: a `Block` type that holds some data, and a way to compute a SHA-256 fingerprint of its contents. You'll also see, hands-on, why the hash has to be computed over a *canonical* byte representation — not just "whatever Go's default struct printing gives you."

## Concepts

### Why hash a block at all?

A hash function takes data of any size and produces a fixed-size fingerprint. Change one byte of the input and the fingerprint comes out completely different (the avalanche effect) — there's no "slightly different" hash for "slightly different" data. That's what makes a block's hash useful as a tamper-detector: if you recompute a block's hash later and it doesn't match what you stored, something changed.

```
   Block{Data: "hello"}   --SHA-256-->   9595c9df90075148eb06860365df33584b75bff...
   Block{Data: "hellp"}   --SHA-256-->   a1e91d1f... (completely different, not "close")
```

### Canonical serialization

Before you can hash a struct, you need to turn it into bytes *deterministically* — the same struct must always produce the same bytes, every time, on every machine. `encoding/gob` (or `encoding/json`) handles this for you as long as you always serialize fields in the same order, which Go's own encoders already guarantee for a fixed struct definition.

### Don't hash your own hash

A block's `Hash` field is *computed from* its other fields — so when computing that hash, you must exclude the `Hash` field itself from what gets serialized, or you'd be trying to hash a value that depends on itself.

## Interface to implement

```go
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

// Hash returns the SHA-256 hash of data.
func Hash(data []byte) []byte

// Serialize returns a deterministic byte representation of the block,
// suitable for hashing. It must NOT include the block's own Hash field.
func (b *Block) Serialize() []byte

// ComputeHash computes and returns this block's hash from its other fields.
// It does not modify b.Hash -- the caller decides when to assign it.
func (b *Block) ComputeHash() []byte

// NewBlock creates a block with the given data and previous hash, and
// sets its Hash field by calling ComputeHash.
func NewBlock(data []byte, prevBlockHash []byte) *Block
```

## Hints

- Use `crypto/sha256` for `Hash` — don't implement SHA-256 yourself.
- `Serialize` can build a byte slice by concatenating `Timestamp` (converted to bytes), `Data`, and `PrevBlockHash` — order just needs to be consistent every time.
- `ComputeHash` should build a temporary copy (or a temporary struct) with `Hash` left empty/nil before serializing, so the current value of `Hash` never leaks into what gets hashed.
- Write a quick manual test: create two blocks with data differing by one character, and print both hashes side by side. Confirm they share no visible structure.

## Run the tests

```bash
cd starter/task-01-hash-a-block
go test ./...
```

All tests must pass before moving to Task 02.
