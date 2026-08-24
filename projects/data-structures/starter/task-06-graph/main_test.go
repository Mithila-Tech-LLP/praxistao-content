package main

import (
	"reflect"
	"sort"
	"testing"
)

// sortedInts returns a sorted copy — used when traversal order among siblings is implementation-defined.
func sortedInts(s []int) []int {
	cp := make([]int, len(s))
	copy(cp, s)
	sort.Ints(cp)
	return cp
}

func TestGraph_AddEdge_Undirected(t *testing.T) {
	g := NewGraph()
	g.AddEdge(1, 2)

	// both directions must be present
	found12 := false
	for _, n := range g.adj[1] {
		if n == 2 {
			found12 = true
		}
	}
	found21 := false
	for _, n := range g.adj[2] {
		if n == 1 {
			found21 = true
		}
	}
	if !found12 {
		t.Error("AddEdge(1,2): 2 not found in adj[1]")
	}
	if !found21 {
		t.Error("AddEdge(1,2): 1 not found in adj[2]")
	}
}

func TestGraph_BFS_SimpleChain(t *testing.T) {
	// 1 - 2 - 3 - 4
	g := NewGraph()
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)

	got := g.BFS(1)
	// BFS from 1 must visit 1 first
	if len(got) == 0 || got[0] != 1 {
		t.Errorf("BFS should start with 1, got %v", got)
	}
	// All 4 nodes must be visited
	if len(got) != 4 {
		t.Errorf("BFS should visit 4 nodes, got %v", got)
	}
	// Sorted result must equal [1,2,3,4]
	if !reflect.DeepEqual(sortedInts(got), []int{1, 2, 3, 4}) {
		t.Errorf("BFS nodes = %v, want all of [1,2,3,4]", got)
	}
}

func TestGraph_BFS_Order(t *testing.T) {
	//   1
	//  / \
	// 2   3
	// |
	// 4
	g := NewGraph()
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)

	got := g.BFS(1)
	if got[0] != 1 {
		t.Errorf("BFS[0] should be 1, got %d", got[0])
	}
	// Nodes 2 and 3 must appear before node 4 (level-order property)
	idx := func(v int) int {
		for i, n := range got {
			if n == v {
				return i
			}
		}
		return -1
	}
	if idx(2) >= idx(4) || idx(3) >= idx(4) {
		t.Errorf("BFS level order violated: 2 and 3 should appear before 4, got %v", got)
	}
}

func TestGraph_DFS_SimpleChain(t *testing.T) {
	// 1 - 2 - 3 - 4
	g := NewGraph()
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)

	got := g.DFS(1)
	if len(got) == 0 || got[0] != 1 {
		t.Errorf("DFS should start with 1, got %v", got)
	}
	if len(got) != 4 {
		t.Errorf("DFS should visit 4 nodes, got %v", got)
	}
	if !reflect.DeepEqual(sortedInts(got), []int{1, 2, 3, 4}) {
		t.Errorf("DFS nodes = %v, want all of [1,2,3,4]", got)
	}
}

func TestGraph_DFS_NoRepeat(t *testing.T) {
	g := NewGraph()
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 1) // cycle

	got := g.DFS(1)
	if len(got) != 3 {
		t.Errorf("DFS on cycle should visit 3 unique nodes, got %v", got)
	}
	seen := map[int]bool{}
	for _, n := range got {
		if seen[n] {
			t.Errorf("DFS visited node %d more than once in %v", n, got)
		}
		seen[n] = true
	}
}

func TestGraph_HasPath_Connected(t *testing.T) {
	g := NewGraph()
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)

	if !g.HasPath(1, 4) {
		t.Error("HasPath(1,4) should be true in chain 1-2-3-4")
	}
	if !g.HasPath(4, 1) {
		t.Error("HasPath(4,1) should be true (undirected)")
	}
	if !g.HasPath(1, 1) {
		t.Error("HasPath(1,1) should be true (same node)")
	}
}

func TestGraph_HasPath_Disconnected(t *testing.T) {
	g := NewGraph()
	g.AddEdge(1, 2)
	g.AddEdge(3, 4) // separate component

	if g.HasPath(1, 3) {
		t.Error("HasPath(1,3) should be false — different components")
	}
	if g.HasPath(2, 4) {
		t.Error("HasPath(2,4) should be false — different components")
	}
}

func TestGraph_SingleNode(t *testing.T) {
	g := NewGraph()
	// manually register a lone node
	g.adj[5] = []int{}

	bfs := g.BFS(5)
	if !reflect.DeepEqual(bfs, []int{5}) {
		t.Errorf("BFS single node = %v, want [5]", bfs)
	}

	dfs := g.DFS(5)
	if !reflect.DeepEqual(dfs, []int{5}) {
		t.Errorf("DFS single node = %v, want [5]", dfs)
	}

	if !g.HasPath(5, 5) {
		t.Error("HasPath(5,5) on single node should be true")
	}
}

func TestGraph_EmptyGraph(t *testing.T) {
	g := NewGraph()
	if g.HasPath(1, 2) {
		t.Error("HasPath on empty graph should be false")
	}
}
