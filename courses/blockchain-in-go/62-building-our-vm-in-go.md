# Chapter 62: Building Our VM in Go

Chapter 61 specified every opcode precisely: its stack effect, its category, whether it reads `Instruction.Arg`. This chapter turns that specification into real, running Go code — the `gochain/vm` package's `Stack`, `VM`, and `Execute()` loop — and tests every single opcode against the behavior Chapter 61 promised.

## Table of Contents

1. [Package Layout](#1-package-layout)
2. [The Stack Type](#2-the-stack-type)
3. [Encoding Helpers](#3-encoding-helpers)
4. [The VM Type and Constructor](#4-the-vm-type-and-constructor)
5. [The Execute Loop](#5-the-execute-loop)
6. [Implementing Every Opcode](#6-implementing-every-opcode)
7. [Testing Every Opcode](#7-testing-every-opcode)
8. [Running a Full Program End to End](#8-running-a-full-program-end-to-end)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Package Layout

Everything in this chapter lives in `gochain/vm`, alongside the `OpCode` constants and `Instruction` struct the shared course contract already fixed:

```
gochain/vm/
├── vm.go        # OpCode, Instruction, Stack, VM, NewVM, Execute
└── vm_test.go   # one test per opcode, plus a full end-to-end program
```

`gochain/vm` will import `gochain/crypto` (for `OpCheckSig`, calling straight into the `Verify` function built in Chapter 13) and, from Chapter 66 onward, `gochain/storage` (for `OpSLoad`/`OpSStore`). Nothing else — the VM package deliberately does not import `network`, `wallet`, or anything that could introduce non-determinism, exactly as Chapter 59 and 60 argued it must not.

---

## 2. The Stack Type

Recall from Chapter 60: a stack is push-and-pop, Last In First Out, and it is the *only* place values live while a GoChain VM program runs. Here is the whole implementation:

```go
package vm

import "errors"

// ErrStackUnderflow is returned whenever an opcode tries to pop from an
// empty stack — for example, a program that starts with OpAdd before
// pushing any operands. A malformed or malicious program should fail
// safely here, not panic and take the whole node down with it.
var ErrStackUnderflow = errors.New("vm: stack underflow")

// Stack is the VM's single working memory: a LIFO list of byte slices.
// Every opcode either pushes a new value, pops one or more existing
// values, or both.
type Stack struct {
	values [][]byte
}

// Push adds a value to the top of the stack. It never fails — there is
// no fixed capacity here (a real production VM would cap stack depth to
// bound memory use, which Chapter 64's gas accounting effectively does
// for us, since every push costs gas).
func (s *Stack) Push(v []byte) {
	s.values = append(s.values, v)
}

// Pop removes and returns the top value. Popping an empty stack is a
// program error, not a Go-level panic — we return it as a normal error
// so Execute can stop cleanly and report exactly what went wrong.
func (s *Stack) Pop() ([]byte, error) {
	if len(s.values) == 0 {
		return nil, ErrStackUnderflow
	}
	top := len(s.values) - 1
	v := s.values[top]
	s.values = s.values[:top] // shrink the slice; v is still valid
	return v, nil
}

// Len reports how many values are currently on the stack. Tests use
// this to assert a program left exactly the right number of values
// behind.
func (s *Stack) Len() int {
	return len(s.values)
}
```

`Push` and `Pop` are exactly the two operations Chapter 60 introduced. `Pop` returning an `error` instead of panicking matters: a smart contract's program is untrusted input from the network's point of view (even if the author intended it correctly), and a single malformed program should never be able to crash the node process that is validating it. Returning a normal Go `error` lets `Execute()` stop and report the problem instead.

---

## 3. Encoding Helpers

Chapter 61 fixed two encoding conventions: numbers as big-endian `uint64`, booleans as single-byte `0x01`/`0x00`. Here are the small helper functions every arithmetic and comparison opcode relies on:

```go
package vm

import "encoding/binary"

var (
	trueVal  = []byte{1}
	falseVal = []byte{0}
)

// boolToBytes converts a Go bool into the canonical stack encoding every
// comparison opcode uses, so OpEqual, OpGreaterThan, and OpCheckSig all
// agree on what "true" and "false" look like as bytes.
func boolToBytes(b bool) []byte {
	if b {
		return trueVal
	}
	return falseVal
}

// isTruthy implements the "falsy" rule from Chapter 61: empty, or every
// byte zero, counts as false. Anything else counts as true. This is
// what OpJumpIfFalse consults to decide whether to jump.
func isTruthy(v []byte) bool {
	if len(v) == 0 {
		return false
	}
	for _, b := range v {
		if b != 0 {
			return true
		}
	}
	return false
}

// uint64ToBytes encodes a number as 8 big-endian bytes, the format every
// OpPush constant representing a number should already be in.
func uint64ToBytes(n uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, n)
	return buf
}

// bytesToUint64 decodes a byte slice as a big-endian number, right-
// aligning anything shorter than 8 bytes (so a program can push a
// shorter constant like []byte{5} and still have it read as 5, not
// some huge garbage number). GoChain VM numbers top out at uint64 —
// a real production VM would use arbitrary-precision integers instead,
// which we call out here rather than pretend our numbers are unbounded.
func bytesToUint64(b []byte) uint64 {
	if len(b) > 8 {
		b = b[len(b)-8:] // truncate to the low 8 bytes rather than error
	}
	padded := make([]byte, 8)
	copy(padded[8-len(b):], b)
	return binary.BigEndian.Uint64(padded)
}
```

`boolToBytes` and `isTruthy` are the two ends of the same convention: any opcode that produces a boolean uses `boolToBytes`, and the one opcode that consumes one as a condition (`OpJumpIfFalse`) uses `isTruthy`. `uint64ToBytes`/`bytesToUint64` are the arithmetic opcodes' only way of turning stack bytes into numbers and back — keeping that conversion in one place means `OpAdd`, `OpSub`, and `OpGreaterThan` all agree on exactly how a number looks as bytes.

---

## 4. The VM Type and Constructor

This is the exact shape fixed by this volume's shared contract — a `Stack`, the program itself, a program counter, and gas tracking:

```go
package vm

// Instruction is one step of a program: an opcode, plus an optional
// operand. OpPush uses Arg as the constant to push; OpJump and
// OpJumpIfFalse use it as a target instruction index. Every other
// opcode leaves Arg nil and gets everything it needs from the stack.
type Instruction struct {
	Op  OpCode
	Arg []byte
}

// VM is one running instance of GoChain's virtual machine: a stack, a
// program to execute, a program counter pointing at the next
// instruction, and gas accounting to bound how much work it can do.
type VM struct {
	stack    *Stack
	program  []Instruction
	pc       int    // index into program of the next instruction to run
	gasUsed  uint64
	gasLimit uint64
}

// NewVM prepares a VM to run program, capped at gasLimit units of gas.
// Execution does not start until Execute is called — constructing a VM
// never itself runs any contract code.
func NewVM(program []Instruction, gasLimit uint64) *VM {
	return &VM{
		stack:    &Stack{},
		program:  program,
		pc:       0,
		gasUsed:  0,
		gasLimit: gasLimit,
	}
}
```

Every field here maps directly onto something Chapter 60 already introduced conceptually: `stack` is the one working memory, `program` is the fixed list of instructions, `pc` is the program counter. `gasUsed`/`gasLimit` exist from the very first version of this struct because gas accounting is fundamental to running untrusted code safely — Chapter 64 gives it real per-opcode costs, but the fields themselves belong here, in the constructor, from day one.

---

## 5. The Execute Loop

This is the loop Chapter 60 described in prose, now as real Go — fetch, dispatch, update, repeat, until `OpHalt`, an error, or (starting properly in Chapter 64) running out of gas:

```go
package vm

import "fmt"

// ErrProgramCounterOutOfRange means the program counter walked past the
// end of the program (or was jumped somewhere invalid) without ever
// hitting OpHalt — for example, a program with no OpHalt at all.
var ErrProgramCounterOutOfRange = errors.New("vm: program counter out of range")

// Execute runs the VM's program from wherever pc currently is (0, for a
// freshly constructed VM) until it hits OpHalt, a runtime error, or (from
// Chapter 64 onward) exhausts its gas limit. It returns nil on a clean
// OpHalt, or the error that stopped it otherwise.
func (vm *VM) Execute() error {
	for {
		if vm.pc < 0 || vm.pc >= len(vm.program) {
			return ErrProgramCounterOutOfRange
		}

		instr := vm.program[vm.pc]

		// Placeholder gas accounting: every instruction costs a flat 1
		// gas for now, just enough to prove the mechanism works end to
		// end. Chapter 64 replaces this single line with a real,
		// per-opcode cost table and explains exactly why each cost is
		// what it is.
		vm.gasUsed++
		if vm.gasUsed > vm.gasLimit {
			return ErrOutOfGas
		}

		switch instr.Op {
		case OpPush:
			vm.execPush(instr)
		case OpDup:
			if err := vm.execDup(); err != nil {
				return err
			}
		case OpPop:
			if err := vm.execPop(); err != nil {
				return err
			}
		case OpAdd:
			if err := vm.execAdd(); err != nil {
				return err
			}
		case OpSub:
			if err := vm.execSub(); err != nil {
				return err
			}
		case OpEqual:
			if err := vm.execEqual(); err != nil {
				return err
			}
		case OpGreaterThan:
			if err := vm.execGreaterThan(); err != nil {
				return err
			}
		case OpJump:
			vm.execJump(instr)
			continue // execJump already set pc; skip the pc++ below
		case OpJumpIfFalse:
			jumped, err := vm.execJumpIfFalse(instr)
			if err != nil {
				return err
			}
			if jumped {
				continue // pc was already set to the jump target
			}
		case OpCheckSig:
			if err := vm.execCheckSig(); err != nil {
				return err
			}
		case OpSLoad:
			if err := vm.execSLoad(); err != nil {
				return err
			}
		case OpSStore:
			if err := vm.execSStore(); err != nil {
				return err
			}
		case OpHalt:
			return nil
		default:
			return fmt.Errorf("vm: unknown opcode %d at pc=%d", instr.Op, vm.pc)
		}

		vm.pc++
	}
}
```

Notice the `continue` statements for `OpJump` and `OpJumpIfFalse` (when it actually jumps): those two cases set `vm.pc` themselves, to the jump target, so the loop must skip the ordinary `vm.pc++` at the bottom — otherwise a jump to instruction 4 would immediately advance to instruction 5, silently skipping the very instruction the jump was aimed at. Every other opcode falls through to `vm.pc++`, advancing one instruction at a time, exactly like Chapter 60's hand traces.

---

## 6. Implementing Every Opcode

Each opcode gets its own small method, keeping the `switch` in `Execute` easy to scan and each opcode's logic easy to test in isolation.

```go
package vm

import "bytes"

func (vm *VM) execPush(instr Instruction) {
	// OpPush never touches existing stack values — it only adds the
	// instruction's operand as a brand-new top-of-stack value.
	vm.stack.Push(instr.Arg)
}

func (vm *VM) execDup() error {
	top, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	// Copy the bytes rather than reusing the same slice twice, so a
	// later opcode that (hypothetically) mutated one copy in place could
	// never silently corrupt the other.
	dup := make([]byte, len(top))
	copy(dup, top)
	vm.stack.Push(top)
	vm.stack.Push(dup)
	return nil
}

func (vm *VM) execPop() error {
	_, err := vm.stack.Pop()
	return err
}

func (vm *VM) execAdd() error {
	b, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	a, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	vm.stack.Push(uint64ToBytes(bytesToUint64(a) + bytesToUint64(b)))
	return nil
}

// ErrNegativeResult means an OpSub would produce a negative number,
// which GoChain VM's unsigned uint64 encoding cannot represent. Failing
// loudly here is much safer for a financial system than silently
// wrapping around to a huge positive number.
var ErrNegativeResult = errors.New("vm: subtraction would underflow")

func (vm *VM) execSub() error {
	b, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	a, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	av, bv := bytesToUint64(a), bytesToUint64(b)
	if bv > av {
		return ErrNegativeResult
	}
	vm.stack.Push(uint64ToBytes(av - bv))
	return nil
}

func (vm *VM) execEqual() error {
	b, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	a, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	vm.stack.Push(boolToBytes(bytes.Equal(a, b)))
	return nil
}

func (vm *VM) execGreaterThan() error {
	b, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	a, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	vm.stack.Push(boolToBytes(bytesToUint64(a) > bytesToUint64(b)))
	return nil
}

func (vm *VM) execJump(instr Instruction) {
	// The jump target lives in Arg, not the stack — set per Chapter 61.
	vm.pc = int(bytesToUint64(instr.Arg))
}

func (vm *VM) execJumpIfFalse(instr Instruction) (jumped bool, err error) {
	cond, err := vm.stack.Pop()
	if err != nil {
		return false, err
	}
	if !isTruthy(cond) {
		vm.pc = int(bytesToUint64(instr.Arg))
		return true, nil
	}
	return false, nil
}
```

`execAdd` and `execSub` both follow the same shape: pop `b` (the top, pushed most recently), pop `a`, compute, push one result — exactly the `( a b -- result )` stack effect Chapter 61 specified. `execSub`'s guard against a negative result is a deliberate design choice, not an oversight: unsigned integers wrapping around silently (`3 - 5` becoming some enormous number instead of an error) is a real, historically damaging class of bug in financial software, and failing loudly here is far safer than any alternative.

`execJump` and `execJumpIfFalse` both read their target from `instr.Arg`, matching Section 5's `continue` handling in `Execute` — the caller (`Execute`) is responsible for knowing not to also increment `pc` afterward.

Now `OpCheckSig`, which calls into `gochain/crypto` — the one opcode that reaches outside the VM package itself, and only into code Volume 2 already built and tested:

```go
package vm

import "github.com/you/gochain/crypto"

func (vm *VM) execCheckSig() error {
	// Pop order matches Chapter 61's table exactly: pubkey is nearest
	// the top, meaning a program pushes data, then signature, then
	// pubkey, in that order, before this opcode runs.
	pubkey, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	signature, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	data, err := vm.stack.Pop()
	if err != nil {
		return err
	}
	// crypto.Verify was built back in Chapter 13; the VM does not
	// implement any cryptography of its own, it only calls out to code
	// that is already written, tested, and trusted.
	ok := crypto.Verify(pubkey, data, signature)
	vm.stack.Push(boolToBytes(ok))
	return nil
}
```

Finally, `OpSLoad` and `OpSStore`. Chapter 66 (the second half of this volume) wires these to a real, persistent, per-contract storage backend built on `gochain/storage`. For now, so the opcodes exist, are testable, and match Chapter 61's stack effects exactly, they behave as a storage of exactly nothing — reads always return an empty value, writes are accepted but not yet remembered anywhere:

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

func (vm *VM) execSStore() error {
	// key is popped first (nearest the top), matching OpSLoad's
	// convention that the storage key is always closest to the top of
	// the stack. Chapter 66 replaces this body with a real write.
	if _, err := vm.stack.Pop(); err != nil { // key
		return err
	}
	if _, err := vm.stack.Pop(); err != nil { // value
		return err
	}
	return nil
}
```

Both already have the exact stack effect Chapter 61 specified — `( key -- value )` and `( value key -- )` — so no code that calls them will need to change once Chapter 66 gives them a real backing store; only their bodies grow.

---

## 7. Testing Every Opcode

Every opcode gets a focused unit test, checking both the resulting stack and — where relevant — that popping too few values fails safely instead of panicking.

```go
package vm

import (
	"testing"
)

func TestOpPush(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(42)},
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.stack.Len() != 1 {
		t.Fatalf("expected 1 value on stack, got %d", v.stack.Len())
	}
}

func TestOpDup(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(7)},
		{Op: OpDup},
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.stack.Len() != 2 {
		t.Fatalf("expected 2 values after dup, got %d", v.stack.Len())
	}
}

func TestOpPop(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(7)},
		{Op: OpPop},
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.stack.Len() != 0 {
		t.Fatalf("expected empty stack after pop, got %d", v.stack.Len())
	}
}

func TestOpAdd(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(5)},
		{Op: OpPush, Arg: uint64ToBytes(3)},
		{Op: OpAdd},
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := v.stack.Pop()
	if bytesToUint64(got) != 8 {
		t.Fatalf("expected 8, got %d", bytesToUint64(got))
	}
}

func TestOpSub(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(10)},
		{Op: OpPush, Arg: uint64ToBytes(4)},
		{Op: OpSub},
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := v.stack.Pop()
	if bytesToUint64(got) != 6 {
		t.Fatalf("expected 6, got %d", bytesToUint64(got))
	}
}

func TestOpSub_Underflow(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(3)},
		{Op: OpPush, Arg: uint64ToBytes(5)},
		{Op: OpSub}, // 3 - 5, would go negative
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != ErrNegativeResult {
		t.Fatalf("expected ErrNegativeResult, got %v", err)
	}
}

func TestOpEqual(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(9)},
		{Op: OpPush, Arg: uint64ToBytes(9)},
		{Op: OpEqual},
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := v.stack.Pop()
	if !isTruthy(got) {
		t.Fatalf("expected true, got false")
	}
}

func TestOpGreaterThan(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(10)},
		{Op: OpPush, Arg: uint64ToBytes(4)},
		{Op: OpGreaterThan}, // is 10 > 4?
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := v.stack.Pop()
	if !isTruthy(got) {
		t.Fatalf("expected true, got false")
	}
}

func TestOpJump(t *testing.T) {
	program := []Instruction{
		{Op: OpJump, Arg: uint64ToBytes(2)}, // skip straight to index 2
		{Op: OpPush, Arg: uint64ToBytes(999)}, // should NEVER run
		{Op: OpPush, Arg: uint64ToBytes(1)},
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.stack.Len() != 1 {
		t.Fatalf("expected exactly 1 value (the skipped push never ran), got %d", v.stack.Len())
	}
}

func TestOpJumpIfFalse(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: falseVal},
		{Op: OpJumpIfFalse, Arg: uint64ToBytes(4)}, // condition is false: jump
		{Op: OpPush, Arg: uint64ToBytes(111)},       // should be SKIPPED
		{Op: OpHalt},
		{Op: OpPush, Arg: uint64ToBytes(222)}, // instruction index 4
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := v.stack.Pop()
	if bytesToUint64(got) != 222 {
		t.Fatalf("expected 222 (jump taken), got %d", bytesToUint64(got))
	}
}

func TestOpCheckSig(t *testing.T) {
	priv, pub := crypto.GenerateKeyPair() // built in Chapter 13
	data := []byte("pay Bob 5 gochips")
	sig := crypto.Sign(priv, data)

	program := []Instruction{
		{Op: OpPush, Arg: data},
		{Op: OpPush, Arg: sig},
		{Op: OpPush, Arg: pub},
		{Op: OpCheckSig},
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := v.stack.Pop()
	if !isTruthy(got) {
		t.Fatalf("expected signature to verify as true")
	}
}

func TestOpSLoad_EmptyByDefault(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: []byte("some-key")},
		{Op: OpSLoad},
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, _ := v.stack.Pop()
	if len(got) != 0 {
		t.Fatalf("expected empty value before Chapter 66 wires real storage, got %v", got)
	}
}

func TestOpSStore_DoesNotError(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(100)}, // value
		{Op: OpPush, Arg: []byte("balance")},  // key
		{Op: OpSStore},
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.stack.Len() != 0 {
		t.Fatalf("OpSStore should leave nothing on the stack, got %d values", v.stack.Len())
	}
}

func TestExecute_StackUnderflow(t *testing.T) {
	program := []Instruction{
		{Op: OpAdd}, // nothing pushed yet — must fail, not panic
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != ErrStackUnderflow {
		t.Fatalf("expected ErrStackUnderflow, got %v", err)
	}
}

func TestExecute_MissingHalt(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(1)},
		// no OpHalt: pc walks off the end of the program
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != ErrProgramCounterOutOfRange {
		t.Fatalf("expected ErrProgramCounterOutOfRange, got %v", err)
	}
}
```

Running `go test ./gochain/vm/...` at this point exercises every opcode from Chapter 61's table at least once, including the two error paths (stack underflow, and a program that never reaches `OpHalt`) that a naive implementation might let panic instead of fail cleanly.

---

## 8. Running a Full Program End to End

To see the whole machine work together the way Chapter 60's hand-trace did, here is `5 + 3 == 8` running for real:

```go
func TestExecute_FullProgram(t *testing.T) {
	// The exact program traced by hand in Chapter 60, Section 5.
	program := []Instruction{
		{Op: OpPush, Arg: uint64ToBytes(5)},
		{Op: OpPush, Arg: uint64ToBytes(3)},
		{Op: OpAdd},
		{Op: OpPush, Arg: uint64ToBytes(8)},
		{Op: OpEqual},
		{Op: OpHalt},
	}
	v := NewVM(program, 1000)
	if err := v.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := v.stack.Pop()
	if err != nil {
		t.Fatalf("expected a result on the stack: %v", err)
	}
	if !isTruthy(got) {
		t.Fatalf("expected true (5 + 3 == 8), got false")
	}
}
```

```
=== RUN   TestExecute_FullProgram
--- PASS: TestExecute_FullProgram (0.00s)
PASS
ok      github.com/you/gochain/vm    0.004s
```

The stack, at every step of this test, holds exactly the values Chapter 60 traced by hand — `Execute()` is not doing anything conceptually new, only carrying out that same trace mechanically, one opcode at a time.

---

## Summary

- `gochain/vm` implements `Stack` (push/pop over a `[][]byte`), the shared `Instruction`/`VM` types, and an `Execute()` loop that fetches, dispatches via `switch`, and updates the stack and program counter.
- Two encoding helpers — `uint64ToBytes`/`bytesToUint64` for numbers, `boolToBytes`/`isTruthy` for booleans — keep every opcode agreeing on how values look as raw bytes.
- Each opcode from Chapter 61 gets its own small method (`execAdd`, `execJump`, and so on), keeping `Execute`'s `switch` easy to read and each opcode independently testable.
- `OpJump` and `OpJumpIfFalse` (when it jumps) set `vm.pc` directly and `continue` the loop, skipping the ordinary `vm.pc++` so the jump target actually runs next.
- `execSub` deliberately errors on a would-be-negative result rather than wrapping around, since silent unsigned-integer wraparound is a dangerous class of bug in financial code.
- `execCheckSig` calls straight into `gochain/crypto.Verify` (built in Chapter 13) — the VM implements no cryptography of its own.
- `OpSLoad`/`OpSStore` are implemented now with the correct stack effects but no real backing store yet; Chapter 66 replaces only their bodies, not their signatures or call sites.
- A full test suite covers every opcode individually, two failure paths (stack underflow, missing `OpHalt`), and the complete `5 + 3 == 8` program traced by hand in Chapter 60.

---

## Exercises

### Easy

1. Why does `Stack.Pop()` return an `error` instead of letting Go panic on an empty stack? What would panicking risk, given that VM programs come from untrusted transactions?
2. Trace `TestOpJump` by hand, instruction by instruction, and explain in your own words why the second instruction (`OpPush 999`) never executes.
3. What is the difference in behavior between `execSLoad`/`execSStore` right now and what Chapter 66 will make them do? What stays exactly the same?

### Medium

4. Add a new test, `TestOpDup_DoesNotAliasMemory`, that pushes a value, calls `OpDup`, and then (using the exported `Stack` API only) confirms the two resulting stack entries are backed by different underlying byte slices, not the same one.
5. `execAdd` and `execSub` both pop `b` before `a`. Write out, in stack-effect notation, what would go wrong (with a concrete numeric example) if `execSub` accidentally popped in the opposite order.
6. `TestExecute_MissingHalt` expects `ErrProgramCounterOutOfRange`. Modify the test so the program contains an `OpJump` to an index far beyond the end of the program (e.g., 500) instead of just omitting `OpHalt`, and confirm the same error still results. Explain why the same bounds check in `Execute` catches both cases.

### Hard

7. `bytesToUint64` truncates any input longer than 8 bytes to its low 8 bytes rather than returning an error. Write a test that pushes a 9-byte constant via `OpPush` and demonstrates exactly what number `execAdd` computes from it. Is this truncation behavior safe for a contract dealing with real value transfers? Explain your reasoning.
8. Implement (in Go, as if adding to this chapter's file) a new opcode `OpNot` with stack effect `( a -- result )` that pushes the boolean opposite of a popped boolean value, using `isTruthy` and `boolToBytes`. Add it to the `OpCode` enum, the `Execute` switch, and write a unit test for it in the same style as this chapter's other opcode tests.
9. The gas accounting in this chapter's `Execute` charges a flat 1 gas per instruction, regardless of opcode. Write a test program that would behave identically under this flat-cost scheme but very differently once Chapter 64's real per-opcode costs are wired in (hint: think about which opcodes Chapter 64 will make expensive). Explain, without writing Chapter 64's code, what you expect to change.
