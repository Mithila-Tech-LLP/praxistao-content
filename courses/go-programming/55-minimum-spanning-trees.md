# Chapter 55: Minimum Spanning Trees — Kruskal's and Prim's

A **minimum spanning tree (MST)** connects all vertices of a weighted undirected graph using exactly V-1 edges with the smallest possible total weight. MSTs appear everywhere: laying fiber optic cables between cities, designing road networks, clustering data points, and building peer-to-peer network overlays.

## Table of Contents

1. [What is a MST?](#1-what-is-a-mst)
2. [Union-Find (Disjoint Set Union)](#2-union-find-disjoint-set-union)
3. [Kruskal's Algorithm](#3-kruskals-algorithm)
4. [Prim's Algorithm](#4-prims-algorithm)
5. [Complexity Comparison](#5-complexity-comparison)
6. [Worked Example: 6-Node Network](#6-worked-example-6-node-network)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. What is a MST?

Given a connected, weighted, undirected graph G = (V, E), a **spanning tree** is a subgraph that:
- Includes all V vertices
- Is a tree (connected, no cycles)
- Has exactly V-1 edges

A **minimum spanning tree** is the spanning tree with the smallest sum of edge weights.

```
Use cases:
  Network cabling    — minimum total wire to connect all offices
  Road layout        — minimum asphalt to connect all towns
  Cluster analysis   — MST-based clustering (cut the k-1 heaviest MST edges → k clusters)
  Approximate TSP    — MST gives a 2-approximation for the travelling salesman problem
  Image segmentation — connect pixels with similar colors at minimum cost
```

**Key property (cut property):** For any cut of the graph into two sets S and V-S, the minimum-weight edge crossing the cut is always in some MST.

---

## 2. Union-Find (Disjoint Set Union)

Kruskal's algorithm needs a fast way to check whether adding an edge would form a cycle. Two vertices form a cycle if they are already in the same component. Union-Find tracks connected components in near-O(1) per operation.

```go
package mst

// UnionFind supports union by rank and path compression.
// Both techniques combined give inverse-Ackermann O(α(n)) ≈ O(1) amortized.
type UnionFind struct {
    parent []int
    rank   []int
    count  int // number of distinct components
}

func NewUnionFind(n int) *UnionFind {
    uf := &UnionFind{
        parent: make([]int, n),
        rank:   make([]int, n),
        count:  n,
    }
    for i := range uf.parent {
        uf.parent[i] = i // each node is its own root
    }
    return uf
}

// Find returns the root of x's component (with path compression).
func (uf *UnionFind) Find(x int) int {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x]) // path compression: flatten tree
    }
    return uf.parent[x]
}

// Union merges the components of x and y.
// Returns true if they were in different components (merge happened).
func (uf *UnionFind) Union(x, y int) bool {
    rx, ry := uf.Find(x), uf.Find(y)
    if rx == ry {
        return false // already same component — adding this edge makes a cycle
    }
    // Union by rank: attach smaller tree under larger tree
    switch {
    case uf.rank[rx] < uf.rank[ry]:
        uf.parent[rx] = ry
    case uf.rank[rx] > uf.rank[ry]:
        uf.parent[ry] = rx
    default:
        uf.parent[ry] = rx
        uf.rank[rx]++
    }
    uf.count--
    return true
}

// Connected reports whether x and y are in the same component.
func (uf *UnionFind) Connected(x, y int) bool {
    return uf.Find(x) == uf.Find(y)
}
```

### Why path compression + union by rank?

```
Without optimizations:
  Find can walk a chain O(n) nodes long → O(n) per operation

With union by rank alone:
  Tree height bounded by O(log n) → O(log n) per Find

With path compression alone:
  Amortized O(log n) per Find

With both:
  Amortized O(α(n)) per operation — α grows so slowly it's effectively O(1)
  For n = 10^80 (atoms in universe), α(n) = 5
```

---

## 3. Kruskal's Algorithm

**Strategy**: sort all edges by weight, then greedily add the cheapest edge that doesn't form a cycle.

```
Algorithm:
  1. Sort all edges by weight (ascending)
  2. Initialize Union-Find with V components
  3. For each edge (u, v, w) in sorted order:
       if u and v are in different components:
           add edge to MST
           union their components
       stop when MST has V-1 edges
```

```go
package mst

import "sort"

// Edge represents a weighted undirected edge.
type Edge struct {
    U, V   int
    Weight int
}

// KruskalMST returns the MST edges and total weight using Kruskal's algorithm.
// n = number of vertices (0-indexed). edges = all undirected edges.
// Returns (nil, -1) if the graph is disconnected.
func KruskalMST(n int, edges []Edge) ([]Edge, int) {
    // Step 1: sort edges by weight
    sorted := make([]Edge, len(edges))
    copy(sorted, edges)
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].Weight < sorted[j].Weight
    })

    uf := NewUnionFind(n)
    var mstEdges []Edge
    totalWeight := 0

    // Step 2: greedily pick edges that don't form cycles
    for _, e := range sorted {
        if len(mstEdges) == n-1 {
            break // MST is complete
        }
        if uf.Union(e.U, e.V) {
            mstEdges = append(mstEdges, e)
            totalWeight += e.Weight
        }
    }

    // A spanning tree needs exactly n-1 edges; fewer means disconnected graph
    if len(mstEdges) < n-1 {
        return nil, -1
    }
    return mstEdges, totalWeight
}
```

### Kruskal trace on a small graph

```
Graph (5 vertices, 7 edges):
  0-1: 4    0-2: 3
  1-2: 1    1-3: 2
  2-3: 4    2-4: 5
  3-4: 7

Sorted edges: (1-2,1), (1-3,2), (0-2,3), (0-1,4), (2-3,4), (2-4,5), (3-4,7)

Components: {0} {1} {2} {3} {4}

Pick (1-2, w=1): union 1 and 2 → {0} {1,2} {3} {4}    MST: [(1-2,1)]
Pick (1-3, w=2): union 1 and 3 → {0} {1,2,3} {4}       MST: [(1-2,1),(1-3,2)]
Pick (0-2, w=3): union 0 and 2 → {0,1,2,3} {4}          MST: [(1-2,1),(1-3,2),(0-2,3)]
Pick (0-1, w=4): Find(0)==Find(1) → SKIP (cycle!)
Pick (2-3, w=4): Find(2)==Find(3) → SKIP (cycle!)
Pick (2-4, w=5): union 2 and 4 → {0,1,2,3,4}            MST: [(1-2,1),(1-3,2),(0-2,3),(2-4,5)]

MST has 4 = n-1 edges. Done. Total weight = 1+2+3+5 = 11
```

---

## 4. Prim's Algorithm

**Strategy**: grow the MST one vertex at a time. Start from any vertex, always add the cheapest edge that connects a new vertex to the current tree.

Prim's is similar to Dijkstra: instead of tracking distance from source, track the minimum edge weight to reach each vertex from the current tree.

```go
package mst

import "container/heap"

// heapItem is an entry in the min-heap used by Prim's algorithm.
type heapItem struct {
    vertex int
    weight int
    from   int // which MST vertex this edge comes from
}

type minHeap []heapItem

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].weight < h[j].weight }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(heapItem)) }
func (h *minHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

// AdjEdge is a neighbor in an adjacency list.
type AdjEdge struct {
    To     int
    Weight int
}

// PrimMST returns the MST edges and total weight using Prim's algorithm.
// graph[u] = list of edges from u. start = source vertex (any vertex works).
// Returns (nil, -1) if the graph is disconnected.
func PrimMST(n int, graph [][]AdjEdge, start int) ([]Edge, int) {
    inMST := make([]bool, n)
    mstEdges := make([]Edge, 0, n-1)
    totalWeight := 0

    h := &minHeap{{vertex: start, weight: 0, from: -1}}
    heap.Init(h)

    for h.Len() > 0 && len(mstEdges) < n-1 {
        item := heap.Pop(h).(heapItem)
        u := item.vertex

        if inMST[u] {
            continue // already in tree — skip stale heap entry
        }
        inMST[u] = true

        if item.from != -1 {
            mstEdges = append(mstEdges, Edge{U: item.from, V: u, Weight: item.weight})
            totalWeight += item.weight
        }

        // Push all neighbors not yet in MST
        for _, e := range graph[u] {
            if !inMST[e.To] {
                heap.Push(h, heapItem{vertex: e.To, weight: e.Weight, from: u})
            }
        }
    }

    if len(mstEdges) < n-1 {
        return nil, -1 // disconnected
    }
    return mstEdges, totalWeight
}

// BuildAdjList converts an edge list to an adjacency list for undirected graph.
func BuildAdjList(n int, edges []Edge) [][]AdjEdge {
    graph := make([][]AdjEdge, n)
    for _, e := range edges {
        graph[e.U] = append(graph[e.U], AdjEdge{To: e.V, Weight: e.Weight})
        graph[e.V] = append(graph[e.V], AdjEdge{To: e.U, Weight: e.Weight})
    }
    return graph
}
```

### Prim trace on the same graph

```
Graph (5 vertices): 0-1:4, 0-2:3, 1-2:1, 1-3:2, 2-3:4, 2-4:5, 3-4:7
Start: vertex 0

Heap: [(0, w=0)]

Pop (0, w=0): add 0 to MST. Push neighbors: (1,w=4,from=0), (2,w=3,from=0)
Heap: [(2,w=3,from=0), (1,w=4,from=0)]

Pop (2, w=3, from=0): add edge 0-2 to MST. Push: (0→in), (1,w=1,from=2), (3,w=4,from=2), (4,w=5,from=2)
Heap: [(1,w=1,from=2), (1,w=4,from=0), (3,w=4,from=2), (4,w=5,from=2)]

Pop (1, w=1, from=2): add edge 2-1 to MST. Push: (0→in), (2→in), (3,w=2,from=1)
Heap: [(3,w=2,from=1), (1,w=4,from=0 — stale), (3,w=4,from=2), (4,w=5,from=2)]

Pop (3, w=2, from=1): add edge 1-3 to MST. Push: (1→in), (2→in), (4,w=7,from=3)
Heap: [(1→stale), (3→stale), (4,w=5,from=2), (4,w=7,from=3)]

Pop (1→stale): inMST[1]=true → skip
Pop (3→stale): inMST[3]=true → skip
Pop (4, w=5, from=2): add edge 2-4 to MST. Done.

MST edges: (0-2,3), (2-1,1), (1-3,2), (2-4,5)  Total = 11 ✓
```

---

## 5. Complexity Comparison

| | Kruskal's | Prim's (binary heap) | Prim's (Fibonacci heap) |
|---|---|---|---|
| **Time** | O(E log E) | O(E log V) | O(E + V log V) |
| **Space** | O(V + E) | O(V + E) | O(V + E) |
| **Best for** | Sparse graphs | Dense graphs | Dense graphs (theory) |
| **Data structure** | Union-Find | Min-heap | Fibonacci heap |
| **Implementation** | Simple | Moderate | Complex |

Since E ≤ V², log E ≤ 2 log V, so O(E log E) = O(E log V). Both algorithms have the same asymptotic complexity for practical purposes.

```
Choosing between them:

Kruskal's wins when:
  - Graph is sparse (E ≈ V)
  - Edges are already sorted (pre-sorted input)
  - You need a simple, easy-to-debug implementation

Prim's wins when:
  - Graph is dense (E ≈ V²) — the heap has at most V entries at a time
  - Graph is given as an adjacency matrix
  - You need to stream: add one vertex at a time (online MST growth)
```

---

## 6. Worked Example: 6-Node Network

A telecom company wants to lay fiber between 6 data centers (0-5). Find the minimum cable to connect all of them.

```
         2
    0 ——————— 1
    |  \      |
  6 |   \ 3  | 5
    |    \   |
    2     3——4
    |      \ |
  8 |    4  \| 7
    5 ————————+
         9
```

```go
package main

import (
    "fmt"
    "mst" // our package above
)

func main() {
    edges := []mst.Edge{
        {U: 0, V: 1, Weight: 2},
        {U: 0, V: 2, Weight: 6},
        {U: 0, V: 3, Weight: 3},
        {U: 1, V: 3, Weight: 5},
        {U: 1, V: 4, Weight: 8},
        {U: 2, V: 5, Weight: 8},
        {U: 3, V: 4, Weight: 4},
        {U: 3, V: 5, Weight: 9},
        {U: 4, V: 5, Weight: 7},
    }
    n := 6

    // Kruskal's
    kEdges, kWeight := mst.KruskalMST(n, edges)
    fmt.Println("Kruskal MST:")
    for _, e := range kEdges {
        fmt.Printf("  %d - %d  (weight %d)\n", e.U, e.V, e.Weight)
    }
    fmt.Println("Total weight:", kWeight)

    // Prim's (from vertex 0)
    graph := mst.BuildAdjList(n, edges)
    pEdges, pWeight := mst.PrimMST(n, graph, 0)
    fmt.Println("\nPrim MST:")
    for _, e := range pEdges {
        fmt.Printf("  %d - %d  (weight %d)\n", e.U, e.V, e.Weight)
    }
    fmt.Println("Total weight:", pWeight)
}
```

Expected output:
```
Kruskal MST:
  0 - 1  (weight 2)
  0 - 3  (weight 3)
  3 - 4  (weight 4)
  1 - 3  (weight 5)  ← skipped: cycle check
  0 - 2  (weight 6)
  4 - 5  (weight 7)
Total weight: 22

Prim MST:
  0 - 1  (weight 2)
  0 - 3  (weight 3)
  3 - 4  (weight 4)
  0 - 2  (weight 6)
  4 - 5  (weight 7)
Total weight: 22
```

Both algorithms produce the same total weight (22), though the edge sets may differ when multiple MSTs exist with the same total weight.

---

## Summary

- An MST connects all V vertices with V-1 edges at minimum total cost.
- **Kruskal's**: sort edges, use Union-Find to skip cycle-forming edges. O(E log E). Best for sparse graphs.
- **Prim's**: grow the tree from a seed vertex using a min-heap. O(E log V). Best for dense graphs.
- **Union-Find** with path compression + union by rank gives near-O(1) per operation. It is the core data structure that makes Kruskal's efficient.
- Both algorithms are greedy and provably correct via the cut property.
- MSTs are not unique when multiple edges share the same weight.

---

## Exercises

### Easy

1. Trace Kruskal's algorithm on the following graph. Show the Union-Find state after each edge is processed:
   ```
   Vertices: 0, 1, 2, 3
   Edges: (0-1, 5), (0-2, 3), (1-2, 4), (1-3, 6), (2-3, 2)
   ```

2. Implement `CountMSTs(n int, edges []Edge) int` that counts how many distinct minimum spanning trees exist. (Hint: edges with equal weights are interchangeable if they connect the same pair of components.)

3. Add a `String() string` method to the `UnionFind` struct that prints the current component assignment in a readable format, e.g., `"[0:root rank=1] [1→0] [2→0] [3:root rank=0]"`. Use it to trace the Kruskal example above step by step.

### Medium

4. Implement **Borůvka's algorithm**: in each phase, every component picks its cheapest outgoing edge and adds it to the MST. Merge components. Repeat until one component remains. This runs in O(E log V) but is naturally parallel — each component can pick independently. Implement it and verify it gives the same MST weight as Kruskal's.

5. **Second minimum spanning tree**: given a graph, find the spanning tree with the second-lowest total weight. One approach: for each edge e not in the MST, the spanning tree formed by adding e and removing the maximum-weight edge on the path between e's endpoints has a higher or equal cost. Find the minimum such cost across all non-MST edges.

6. **Online MST with edge deletions**: given an MST and a deleted edge (u, v), efficiently repair the MST. The MST splits into two components; find the minimum-weight edge that reconnects them. Implement this with a BFS/DFS over the original graph. What is the time complexity of a single repair?

### Hard

7. Implement **Prim's algorithm with a Fibonacci heap** (or simulate it with a decrease-key lazy heap). The key operation is `decrease-key`: when a shorter edge to vertex v is found, update its priority in the heap. A simple binary heap discards stale entries (the approach above); a proper Fibonacci heap achieves O(1) amortized decrease-key. Benchmark both on a dense random graph with 1000 vertices.

8. **MST-based clustering**: implement k-means-style clustering using an MST. Build the MST of a set of 2D points (where edge weight is Euclidean distance), then remove the k-1 heaviest edges to produce k clusters. Write a program that reads points from stdin, builds the MST with Kruskal's, and outputs cluster assignments. Test it on two clearly separated Gaussian clusters and verify the algorithm separates them cleanly.
