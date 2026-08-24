package main

import "testing"

func lastOf(s []int64) int64 {
	if len(s) == 0 {
		return -999999
	}
	return s[len(s)-1]
}

func TestSimpleAddition(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: 3},
		{Op: OpPush, Arg: 4},
		{Op: OpAdd},
		{Op: OpHalt},
	}
	vm := NewVM(program, 1000)
	if err := vm.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := lastOf(vm.Stack()); got != 7 {
		t.Fatalf("expected 3+4=7, got %d", got)
	}
}

func TestComparisonProgram(t *testing.T) {
	// PUSH 3, PUSH 4, ADD, PUSH 2, GT, HALT  -> (3+4) > 2 -> true (1)
	program := []Instruction{
		{Op: OpPush, Arg: 3},
		{Op: OpPush, Arg: 4},
		{Op: OpAdd},
		{Op: OpPush, Arg: 2},
		{Op: OpGreaterThan},
		{Op: OpHalt},
	}
	vm := NewVM(program, 1000)
	if err := vm.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := lastOf(vm.Stack()); got != 1 {
		t.Fatalf("expected stack top to be 1 (true), got %d", got)
	}
}

func TestSubtractionAndEquality(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: 10},
		{Op: OpPush, Arg: 3},
		{Op: OpSub}, // 10 - 3 = 7
		{Op: OpPush, Arg: 7},
		{Op: OpEqual},
		{Op: OpHalt},
	}
	vm := NewVM(program, 1000)
	if err := vm.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := lastOf(vm.Stack()); got != 1 {
		t.Fatalf("expected 10-3 == 7 to be true, got %d", got)
	}
}

func TestJumpIfFalseSkipsBranch(t *testing.T) {
	// if (5 == 6) { push 111 } else { push 222 }; halt
	program := []Instruction{
		{Op: OpPush, Arg: 5},
		{Op: OpPush, Arg: 6},
		{Op: OpEqual},               // [2]
		{Op: OpJumpIfFalse, Arg: 6}, // [3] if false, jump to index 6
		{Op: OpPush, Arg: 111},      // [4] "then" branch
		{Op: OpJump, Arg: 7},        // [5] skip the else branch
		{Op: OpPush, Arg: 222},      // [6] "else" branch
		{Op: OpHalt},                // [7]
	}
	vm := NewVM(program, 1000)
	if err := vm.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := lastOf(vm.Stack()); got != 222 {
		t.Fatalf("expected the else branch (222) to run since 5 != 6, got %d", got)
	}
}

func TestOutOfGasOnInfiniteLoop(t *testing.T) {
	// An infinite loop: jump to itself forever.
	program := []Instruction{
		{Op: OpJump, Arg: 0},
	}
	vm := NewVM(program, 50) // small, finite gas budget
	err := vm.Execute()
	if err != ErrOutOfGas {
		t.Fatalf("expected ErrOutOfGas, got %v", err)
	}
}

func TestEnoughGasForNormalProgram(t *testing.T) {
	program := []Instruction{
		{Op: OpPush, Arg: 1},
		{Op: OpPush, Arg: 2},
		{Op: OpAdd},
		{Op: OpHalt},
	}
	// Exactly 4 instructions -- must succeed with gas >= 4.
	vm := NewVM(program, 4)
	if err := vm.Execute(); err != nil {
		t.Fatalf("expected a short program to run within a small gas budget, got: %v", err)
	}
}

func TestStackUnderflowIsAnError(t *testing.T) {
	program := []Instruction{
		{Op: OpAdd}, // nothing pushed yet!
	}
	vm := NewVM(program, 1000)
	if err := vm.Execute(); err == nil {
		t.Fatal("expected an error popping from an empty stack")
	}
}
