package main

import (
	"reflect"
	"testing"
)

func TestBST_EmptyTree(t *testing.T) {
	tree := &BST{}

	if tree.Search(5) {
		t.Error("Search on empty tree should return false")
	}

	inorder := tree.InOrder()
	if len(inorder) != 0 {
		t.Errorf("InOrder on empty tree should return empty slice, got %v", inorder)
	}

	_, ok := tree.Min()
	if ok {
		t.Error("Min on empty tree should return false")
	}

	_, ok = tree.Max()
	if ok {
		t.Error("Max on empty tree should return false")
	}
}

func TestBST_SingleNode(t *testing.T) {
	tree := &BST{}
	tree.Insert(42)

	if !tree.Search(42) {
		t.Error("Search should find 42 after inserting it")
	}
	if tree.Search(0) {
		t.Error("Search for 0 should return false")
	}

	inorder := tree.InOrder()
	if !reflect.DeepEqual(inorder, []int{42}) {
		t.Errorf("InOrder should return [42], got %v", inorder)
	}

	min, ok := tree.Min()
	if !ok || min != 42 {
		t.Errorf("Min should return (42, true), got (%d, %v)", min, ok)
	}

	max, ok := tree.Max()
	if !ok || max != 42 {
		t.Errorf("Max should return (42, true), got (%d, %v)", max, ok)
	}
}

func TestBST_Insert_And_Search(t *testing.T) {
	tree := &BST{}
	values := []int{5, 3, 7, 1, 4, 6, 8}
	for _, v := range values {
		tree.Insert(v)
	}

	for _, v := range values {
		if !tree.Search(v) {
			t.Errorf("Search should find %d after inserting it", v)
		}
	}

	notInserted := []int{0, 2, 9, 100}
	for _, v := range notInserted {
		if tree.Search(v) {
			t.Errorf("Search should not find %d (not inserted)", v)
		}
	}
}

func TestBST_InOrder_Sorted(t *testing.T) {
	tree := &BST{}
	for _, v := range []int{5, 3, 7, 1, 4, 6, 8} {
		tree.Insert(v)
	}

	got := tree.InOrder()
	want := []int{1, 3, 4, 5, 6, 7, 8}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder = %v, want %v", got, want)
	}
}

func TestBST_InOrder_SortedInsertion(t *testing.T) {
	// inserting already-sorted values creates a skewed tree
	tree := &BST{}
	for _, v := range []int{1, 2, 3, 4, 5} {
		tree.Insert(v)
	}

	got := tree.InOrder()
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder (sorted insert) = %v, want %v", got, want)
	}
}

func TestBST_InOrder_ReverseInsertion(t *testing.T) {
	tree := &BST{}
	for _, v := range []int{5, 4, 3, 2, 1} {
		tree.Insert(v)
	}

	got := tree.InOrder()
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder (reverse insert) = %v, want %v", got, want)
	}
}

func TestBST_Duplicates(t *testing.T) {
	tree := &BST{}
	tree.Insert(5)
	tree.Insert(5) // duplicate

	// InOrder must NOT contain a duplicate
	got := tree.InOrder()
	if len(got) != 1 {
		t.Errorf("InOrder after inserting duplicate should have 1 element, got %v", got)
	}
	if got[0] != 5 {
		t.Errorf("InOrder after duplicate insert should be [5], got %v", got)
	}
}

func TestBST_MinMax(t *testing.T) {
	tree := &BST{}
	for _, v := range []int{5, 3, 7, 1, 9, 2, 8} {
		tree.Insert(v)
	}

	min, ok := tree.Min()
	if !ok || min != 1 {
		t.Errorf("Min should return (1, true), got (%d, %v)", min, ok)
	}

	max, ok := tree.Max()
	if !ok || max != 9 {
		t.Errorf("Max should return (9, true), got (%d, %v)", max, ok)
	}
}
