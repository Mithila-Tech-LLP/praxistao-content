# Task 02: Link the Chain

## What you will build

A `Blockchain` type that links blocks together by hash, starting from a special first block (the genesis block), plus validation logic that detects tampering anywhere in the chain's history.

## Concepts

### Linking by hash

A blockchain is a chain the moment each block stores the hash of the block before it. That's the entire mechanism:

```
Block 0 (genesis)         Block 1                    Block 2
+----------------+        +----------------+         +----------------+
| PrevHash: 0000  |        | PrevHash: --*  |         | PrevHash: --*  |
| Hash: aaa1      |<-------| Hash: bbb2     |<--------| Hash: ccc3     |
+----------------+        +----------------+         +----------------+
```

The genesis block has no predecessor, so its `PrevBlockHash` is conventionally all zero bytes.

### Validation means recomputing, not trusting

`ValidateBlock` should never just check "does this look right?" — it should independently recompute the block's hash from its actual contents and compare it to the stored `Hash` field, and separately check that `PrevBlockHash` really does match the previous block's real, stored hash. If you tamper with any field of an old block, its recomputed hash changes, which breaks the very next block's `PrevBlockHash` check — and every block after that.

## Interface to implement

```go
type Blockchain struct {
	Blocks []*Block
}

// NewGenesisBlock creates the first block in a chain, with PrevBlockHash
// set to a slice of 32 zero bytes.
func NewGenesisBlock() *Block

// NewBlockchain creates a new chain containing only a genesis block.
func NewBlockchain() *Blockchain

// AddBlock creates a new block containing data, links it to the current
// last block, and appends it to the chain.
func (bc *Blockchain) AddBlock(data string)

// IsValid walks the entire chain and returns false if:
//   - any block's stored Hash doesn't match its recomputed hash, or
//   - any block's PrevBlockHash doesn't match the previous block's Hash
func (bc *Blockchain) IsValid() bool
```

## Hints

- Reuse `Block`, `NewBlock`, `Hash`, and `ComputeHash` from Task 01 (copy them into this task's `main.go` and extend as needed).
- `AddBlock` should look at `bc.Blocks[len(bc.Blocks)-1]` to find the current tip's hash.
- For `IsValid`, loop from index 1 onward (skip genesis, which has no predecessor to check against) and compare `bc.Blocks[i].PrevBlockHash` to `bc.Blocks[i-1].Hash`.
- To test tamper-detection: build a valid 3-block chain, confirm `IsValid()` returns `true`, then directly mutate `bc.Blocks[1].Data`, and confirm `IsValid()` now returns `false` — *without* recomputing `bc.Blocks[1].Hash` yourself (that's the whole point: the stored hash and the real contents now disagree).

## Run the tests

```bash
cd starter/task-02-link-the-chain
go test ./...
```

All tests must pass before moving to Task 03.
