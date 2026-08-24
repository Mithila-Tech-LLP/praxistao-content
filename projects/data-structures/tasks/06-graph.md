---
title: Graph
step: 6
difficulty: medium
estimated: 50 min
---

## What You Are Building

A graph is a set of nodes (vertices) connected by edges. Unlike trees, graphs can have cycles, multiple paths between nodes, and disconnected components. You'll build an undirected graph with an adjacency list and implement BFS and DFS traversals.

```
0 — 1 — 3
|   |
2   4
```

## Key Concepts

**Adjacency list** — Represent a graph as a map from each node to its list of neighbours. This is memory-efficient for sparse graphs.

```go
adj = map[int][]int{
    0: {1, 2},
    1: {0, 3, 4},
    2: {0},
    3: {1},
    4: {1},
}
```

**BFS (Breadth-First Search)** — Use a queue. Visit a node, enqueue its unvisited neighbours, repeat. BFS explores level by level — closest nodes first. Always finds the shortest path in an unweighted graph.

```
Start at 0 → visit 0 → enqueue [1, 2]
Visit 1 → enqueue [3, 4] → result: [0, 1, 2, 3, 4]
```

**DFS (Depth-First Search)** — Use a stack (or recursion). Visit a node, then immediately recurse into its first unvisited neighbour, going as deep as possible before backtracking.

**Visited tracking** — Both BFS and DFS require a `visited` set (use `map[int]bool`) to avoid revisiting nodes and getting stuck in cycles.

**Undirected graph** — When you add edge `(A, B)`, add B to A's adjacency list AND A to B's adjacency list.

## Struct Signature

```go
type Graph struct {
    adj map[int][]int
}

func NewGraph() *Graph {
    return &Graph{adj: make(map[int][]int)}
}
```

## Methods to Implement

| Method | Description |
|--------|-------------|
| `AddEdge(from, to int)` | Add an undirected edge |
| `BFS(start int) []int` | Return nodes in BFS visit order |
| `DFS(start int) []int` | Return nodes in DFS visit order |
| `HasPath(from, to int) bool` | True if a path exists between the two nodes |

## Edge Cases to Handle

- `BFS`/`DFS` from a node with no edges: return `[]int{start}`
- `HasPath(x, x)`: return `true` (node reaches itself)
- `HasPath` for disconnected nodes: return `false`
- Graphs with cycles must not loop infinitely

## Example

```go
g := NewGraph()
g.AddEdge(0, 1)
g.AddEdge(0, 2)
g.AddEdge(1, 3)
g.AddEdge(1, 4)

fmt.Println(g.BFS(0)) // [0 1 2 3 4] (or similar BFS order)
fmt.Println(g.DFS(0)) // [0 1 3 4 2] (or similar DFS order)

fmt.Println(g.HasPath(0, 4)) // true
fmt.Println(g.HasPath(3, 2)) // true
fmt.Println(g.HasPath(3, 5)) // false (5 not in graph)
```

## Hints

- **BFS queue**: maintain a `[]int` slice. Dequeue from index 0, enqueue with `append`.
- **DFS**: write a recursive helper `func dfs(node int, visited map[int]bool, result *[]int)`.
- For deterministic test output, sort the neighbour lists before iterating. `sort.Ints(g.adj[node])` inside BFS/DFS gives consistent ordering.
- `HasPath` can reuse BFS or DFS: traverse from `from` and check if `to` appears in the visited set.
