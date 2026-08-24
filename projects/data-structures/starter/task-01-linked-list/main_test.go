package main

import (
	"reflect"
	"testing"
)

func TestLinkedList_AppendAndToSlice(t *testing.T) {
	l := &LinkedList{}
	l.Append(1)
	l.Append(2)
	l.Append(3)
	want := []int{1, 2, 3}
	got := l.ToSlice()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Append: got %v, want %v", got, want)
	}
}

func TestLinkedList_Prepend(t *testing.T) {
	l := &LinkedList{}
	l.Append(2)
	l.Append(3)
	l.Prepend(1)
	l.Prepend(0)
	want := []int{0, 1, 2, 3}
	got := l.ToSlice()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Prepend: got %v, want %v", got, want)
	}
}

func TestLinkedList_ToSliceEmpty(t *testing.T) {
	l := &LinkedList{}
	got := l.ToSlice()
	if len(got) != 0 {
		t.Errorf("ToSlice on empty list: got %v, want []", got)
	}
}

func TestLinkedList_Contains(t *testing.T) {
	l := &LinkedList{}
	l.Append(10)
	l.Append(20)
	l.Append(30)

	if !l.Contains(20) {
		t.Error("Contains(20): expected true")
	}
	if l.Contains(99) {
		t.Error("Contains(99): expected false")
	}
}

func TestLinkedList_ContainsEmpty(t *testing.T) {
	l := &LinkedList{}
	if l.Contains(1) {
		t.Error("Contains on empty list: expected false")
	}
}

func TestLinkedList_DeleteMiddle(t *testing.T) {
	l := &LinkedList{}
	l.Append(1)
	l.Append(2)
	l.Append(3)
	l.Delete(2)
	want := []int{1, 3}
	got := l.ToSlice()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Delete middle: got %v, want %v", got, want)
	}
}

func TestLinkedList_DeleteHead(t *testing.T) {
	l := &LinkedList{}
	l.Append(1)
	l.Append(2)
	l.Append(3)
	l.Delete(1)
	want := []int{2, 3}
	got := l.ToSlice()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Delete head: got %v, want %v", got, want)
	}
}

func TestLinkedList_DeleteTail(t *testing.T) {
	l := &LinkedList{}
	l.Append(1)
	l.Append(2)
	l.Append(3)
	l.Delete(3)
	want := []int{1, 2}
	got := l.ToSlice()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Delete tail: got %v, want %v", got, want)
	}
}

func TestLinkedList_DeleteOnlyElement(t *testing.T) {
	l := &LinkedList{}
	l.Append(42)
	l.Delete(42)
	got := l.ToSlice()
	if len(got) != 0 {
		t.Errorf("Delete only element: got %v, want []", got)
	}
}

func TestLinkedList_DeleteNonExistent(t *testing.T) {
	l := &LinkedList{}
	l.Append(1)
	l.Append(2)
	l.Delete(99) // should not panic or change list
	want := []int{1, 2}
	got := l.ToSlice()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Delete non-existent: got %v, want %v", got, want)
	}
}

func TestLinkedList_DeleteEmpty(t *testing.T) {
	l := &LinkedList{}
	l.Delete(1) // must not panic
}

func TestLinkedList_DeleteFirstOccurrence(t *testing.T) {
	l := &LinkedList{}
	l.Append(5)
	l.Append(5)
	l.Append(5)
	l.Delete(5) // only removes first
	want := []int{5, 5}
	got := l.ToSlice()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Delete first occurrence: got %v, want %v", got, want)
	}
}

func TestLinkedList_MixedOperations(t *testing.T) {
	l := &LinkedList{}
	l.Append(3)
	l.Prepend(1)
	l.Append(5)
	l.Prepend(0)
	// list: 0 -> 1 -> 3 -> 5
	l.Delete(3)
	// list: 0 -> 1 -> 5
	want := []int{0, 1, 5}
	got := l.ToSlice()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Mixed ops: got %v, want %v", got, want)
	}
	if !l.Contains(5) {
		t.Error("Contains(5) should be true")
	}
	if l.Contains(3) {
		t.Error("Contains(3) should be false after delete")
	}
}
