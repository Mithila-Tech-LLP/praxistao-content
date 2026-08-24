# Task 09: A Tiny Contract VM

## What you will build

A minimal stack-based virtual machine — the same core idea behind Bitcoin Script and, at a conceptual level, the Ethereum Virtual Machine — with a small opcode set and gas metering to guarantee every program eventually stops.

## Concepts

### A stack machine

Instructions push and pop values from a single shared stack. `OpPush 5` puts `5` on top of the stack. `OpAdd` pops the top two values and pushes their sum. A whole program is just a sequence of these tiny steps:

```
program: PUSH 3, PUSH 4, ADD, PUSH 2, GT, HALT

stack after PUSH 3:  [3]
stack after PUSH 4:  [3, 4]
stack after ADD:     [7]
stack after PUSH 2:  [7, 2]
stack after GT:      [1]        (7 > 2 is true)
HALT: program stops, result is [1]
```

### Gas: making "runs forever" impossible

A program with a loop that never terminates would hang whatever runs it, forever. Gas fixes this: every instruction costs a small, fixed amount, the caller supplies a maximum gas budget up front, and execution stops the instant that budget is exceeded — turning "infinite loop" into "predictable, safe failure."

## Interface to implement

```go
type OpCode byte

const (
	OpPush OpCode = iota
	OpAdd
	OpSub
	OpEqual
	OpGreaterThan
	OpJump        // unconditional jump to an instruction index
	OpJumpIfFalse // pop one value; jump if it's zero
	OpHalt
)

type Instruction struct {
	Op  OpCode
	Arg int64 // operand for OpPush (the value) or OpJump/OpJumpIfFalse (the target index)
}

var ErrOutOfGas = errors.New("vm: out of gas")

type VM struct {
	// unexported fields
}

// NewVM creates a VM ready to run program with the given gas limit.
func NewVM(program []Instruction, gasLimit int64) *VM

// Execute runs the program until OpHalt, an error, or running out of
// gas (returning ErrOutOfGas in that case).
func (vm *VM) Execute() error

// Stack returns the VM's current stack contents (top of stack last),
// useful for inspecting the result after Execute returns.
func (vm *VM) Stack() []int64
```

## Hints

- Give each opcode a small fixed gas cost (e.g. 1 for arithmetic/comparison, 1 for push, 1 for jumps) and decrement a remaining-gas counter before executing each instruction; return `ErrOutOfGas` the instant it would go negative.
- A `switch` statement over `instruction.Op`, driven by a program counter (`pc`) you increment after each instruction (except jumps, which set `pc` directly), is the cleanest way to implement `Execute`.
- Represent stack values as `int64` for simplicity — no need for byte slices here.
- Write a test for a deliberately infinite loop (`OpJump` back to its own instruction index, forever) with a small gas limit, and confirm `Execute` returns `ErrOutOfGas` rather than hanging your test suite.
- Write a test tracing the example program in the Concepts section above by hand, and confirm `Stack()` returns `[1]` after `Execute`.

## Run the tests

```bash
cd starter/task-09-a-tiny-contract-vm
go test ./...
```

All tests must pass before moving to Task 10.
