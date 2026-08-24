package main

type Graph struct {
	adj map[int][]int
}

func NewGraph() *Graph {
	return &Graph{adj: make(map[int][]int)}
}

func (g *Graph) AddEdge(from, to int) {
	// TODO: undirected — add both directions
}

func (g *Graph) BFS(start int) []int {
	// TODO: breadth-first traversal, return nodes in visit order
	return []int{}
}

func (g *Graph) DFS(start int) []int {
	// TODO: depth-first traversal, return nodes in visit order
	return []int{}
}

func (g *Graph) HasPath(from, to int) bool {
	// TODO: return true if there is any path from 'from' to 'to'
	return false
}
