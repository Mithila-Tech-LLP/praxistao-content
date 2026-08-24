package main

import "errors"

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

const gasPerInstruction = 1

type VM struct {
	program []Instruction
	pc      int
	stack   []int64
	gasLeft int64
}

// NewVM creates a VM ready to run program with the given gas limit.
func NewVM(program []Instruction, gasLimit int64) *VM {
	return &VM{program: program, gasLeft: gasLimit}
}

// Execute runs the program until OpHalt, an error, or running out of
// gas (returning ErrOutOfGas in that case).
func (vm *VM) Execute() error {
	panic("TODO: implement Execute -- a pc-driven loop with a switch over vm.program[vm.pc].Op")
}

// Stack returns the VM's current stack contents (top of stack last).
func (vm *VM) Stack() []int64 {
	return vm.stack
}
