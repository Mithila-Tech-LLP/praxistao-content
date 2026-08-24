# Chapter 39: Graphs

Graphs model relationships: networks, maps, dependencies, social connections. A graph is a set of **vertices** (nodes) connected by **edges**. Unlike trees (acyclic, hierarchical), graphs can have cycles and arbitrary connections. Graph algorithms — BFS, DFS, shortest path, topological sort — appear constantly in backend systems: dependency resolution, route finding, social network analysis, scheduling.

## Table of Contents

1. [Graph Representations](#1-graph-representations)
2. [Graph Traversal — BFS](#2-graph-traversal--bfs)
3. [Graph Traversal — DFS](#3-graph-traversal--dfs)
4. [Cycle Detection](#4-cycle-detection)
5. [Topological Sort](#5-topological-sort)
6. [Shortest Path — BFS and Dijkstra](#6-shortest-path--bfs-and-dijkstra)
7. [Union-Find (Disjoint Set)](#7-union-find-disjoint-set)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Graph Representations

**Graph types:**
```
Undirected: edges have no direction — A-B means A connects to B and B to A
Directed:   edges have direction — A→B means A points to B (not necessarily B to A)
Weighted:   edges have costs — A--(5)--B means the connection costs 5
Unweighted: all edges cost 1 (or just "connected")
```

**Adjacency List** — most common; O(V+E) space:
```go
// Vertices are ints (0 to V-1); edges are neighbors:
type Graph struct {
    vertices int
    adj      [][]int  // adj[v] = list of neighbors of v
}

func NewGraph(v int) *Graph {
    adj := make([][]int, v)
    for i := range adj {
        adj[i] = []int{}
    }
    return &Graph{vertices: v, adj: adj}
}

// Undirected: add edge in both directions
func (g *Graph) AddEdge(u, v int) {
    g.adj[u] = append(g.adj[u], v)
    g.adj[v] = append(g.adj[v], u)
}

// Directed: add edge only one direction
func (g *Graph) AddDirectedEdge(from, to int) {
    g.adj[from] = append(g.adj[from], to)
}
```

**Weighted graph:**
```go
type Edge struct {
    To     int
    Weight int
}

type WeightedGraph struct {
    vertices int
    adj      [][]Edge
}

func NewWeightedGraph(v int) *WeightedGraph {
    adj := make([][]Edge, v)
    return &WeightedGraph{vertices: v, adj: adj}
}

func (g *WeightedGraph) AddEdge(from, to, weight int) {
    g.adj[from] = append(g.adj[from], Edge{to, weight})
    g.adj[to] = append(g.adj[to], Edge{from, weight})  // Undirected
}
```

**Adjacency Matrix** — O(V²) space; fast edge lookup O(1):
```go
type MatrixGraph struct {
    matrix [][]int  // matrix[u][v] = weight (0 = no edge)
}

func NewMatrixGraph(v int) *MatrixGraph {
    m := make([][]int, v)
    for i := range m {
        m[i] = make([]int, v)
    }
    return &MatrixGraph{matrix: m}
}

func (g *MatrixGraph) AddEdge(u, v, weight int) {
    g.matrix[u][v] = weight
    g.matrix[v][u] = weight  // Undirected
}

func (g *MatrixGraph) HasEdge(u, v int) bool { return g.matrix[u][v] != 0 }
```

**String-keyed graph (practical for real systems):**
```go
type StringGraph struct {
    adj map[string][]string
}

func NewStringGraph() *StringGraph {
    return &StringGraph{adj: make(map[string][]string)}
}

func (g *StringGraph) AddEdge(from, to string) {
    g.adj[from] = append(g.adj[from], to)
    g.adj[to] = append(g.adj[to], from)
}

func (g *StringGraph) Neighbors(v string) []string {
    return g.adj[v]
}
```

### Quick Check
> 1. What is the space complexity of an adjacency list vs adjacency matrix?
> 2. When would you prefer an adjacency matrix over an adjacency list?
> 3. In a directed graph, if there's an edge A→B, is there necessarily an edge B→A?

---

## 2. Graph Traversal — BFS

Breadth-first search explores layer by layer — first all neighbors at distance 1, then distance 2, etc. Finds the **shortest path** in unweighted graphs.

```go
// BFS returns all vertices reachable from start, in BFS order.
func (g *Graph) BFS(start int) []int {
    visited := make([]bool, g.vertices)
    var result []int

    queue := []int{start}
    visited[start] = true

    for len(queue) > 0 {
        v := queue[0]
        queue = queue[1:]
        result = append(result, v)

        for _, neighbor := range g.adj[v] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
    return result
}
```

**BFS shortest path:**
```go
// ShortestPath returns the shortest path from start to end (by hop count).
func (g *Graph) ShortestPath(start, end int) []int {
    if start == end {
        return []int{start}
    }

    visited := make([]bool, g.vertices)
    parent := make([]int, g.vertices)
    for i := range parent { parent[i] = -1 }

    queue := []int{start}
    visited[start] = true

    for len(queue) > 0 {
        v := queue[0]
        queue = queue[1:]

        for _, neighbor := range g.adj[v] {
            if !visited[neighbor] {
                visited[neighbor] = true
                parent[neighbor] = v
                if neighbor == end {
                    return reconstructPath(parent, start, end)
                }
                queue = append(queue, neighbor)
            }
        }
    }
    return nil  // No path exists
}

func reconstructPath(parent []int, start, end int) []int {
    var path []int
    for v := end; v != -1; v = parent[v] {
        path = append(path, v)
    }
    // Reverse:
    for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
        path[i], path[j] = path[j], path[i]
    }
    return path
}
```

**BFS on grids (common in interview problems):**
```go
// Shortest path in a 2D grid from (0,0) to (rows-1, cols-1).
// 0 = passable, 1 = wall.
func ShortestGridPath(grid [][]int) int {
    rows, cols := len(grid), len(grid[0])
    if grid[0][0] == 1 || grid[rows-1][cols-1] == 1 {
        return -1
    }

    type Point struct{ r, c int }
    visited := make([][]bool, rows)
    for i := range visited { visited[i] = make([]bool, cols) }

    queue := []Point{{0, 0}}
    visited[0][0] = true
    dist := 1

    // 8 directions — this variant allows diagonal moves (use only the first 4 for up/down/left/right):
    dirs := []Point{{-1,0},{1,0},{0,-1},{0,1},{-1,-1},{-1,1},{1,-1},{1,1}}

    for len(queue) > 0 {
        size := len(queue)
        for i := 0; i < size; i++ {
            p := queue[i]
            if p.r == rows-1 && p.c == cols-1 {
                return dist
            }
            for _, d := range dirs {
                nr, nc := p.r+d.r, p.c+d.c
                if nr >= 0 && nr < rows && nc >= 0 && nc < cols &&
                    !visited[nr][nc] && grid[nr][nc] == 0 {
                    visited[nr][nc] = true
                    queue = append(queue, Point{nr, nc})
                }
            }
        }
        queue = queue[size:]
        dist++
    }
    return -1
}
```

### Quick Check
> 1. Why does BFS find the shortest path in an unweighted graph?
> 2. What data structure does BFS use?
> 3. What is the time complexity of BFS on a graph with V vertices and E edges?

---

## 3. Graph Traversal — DFS

Depth-first search explores as far as possible before backtracking. Used for: cycle detection, topological sort, connected components, path existence.

```go
// DFS recursive — visits all reachable vertices from start:
func (g *Graph) DFS(start int) []int {
    visited := make([]bool, g.vertices)
    var result []int
    g.dfsHelper(start, visited, &result)
    return result
}

func (g *Graph) dfsHelper(v int, visited []bool, result *[]int) {
    visited[v] = true
    *result = append(*result, v)

    for _, neighbor := range g.adj[v] {
        if !visited[neighbor] {
            g.dfsHelper(neighbor, visited, result)
        }
    }
}

// DFS iterative (avoids stack overflow on large graphs):
func (g *Graph) DFSIterative(start int) []int {
    visited := make([]bool, g.vertices)
    var result []int
    stack := []int{start}

    for len(stack) > 0 {
        v := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if visited[v] {
            continue
        }
        visited[v] = true
        result = append(result, v)

        // Push neighbors in reverse order so first neighbor is processed first:
        for i := len(g.adj[v]) - 1; i >= 0; i-- {
            if !visited[g.adj[v][i]] {
                stack = append(stack, g.adj[v][i])
            }
        }
    }
    return result
}
```

**Connected components:**
```go
func (g *Graph) ConnectedComponents() [][]int {
    visited := make([]bool, g.vertices)
    var components [][]int

    for v := 0; v < g.vertices; v++ {
        if !visited[v] {
            var component []int
            g.dfsHelper(v, visited, &component)
            components = append(components, component)
        }
    }
    return components
}
```

### Quick Check
> 1. What data structure does iterative DFS use?
> 2. What is a connected component?
> 3. Time complexity of DFS?

---

## 4. Cycle Detection

**Undirected graph — DFS with parent tracking:**
```go
func (g *Graph) HasCycle() bool {
    visited := make([]bool, g.vertices)
    for v := 0; v < g.vertices; v++ {
        if !visited[v] {
            if g.dfsCycle(v, -1, visited) {
                return true
            }
        }
    }
    return false
}

func (g *Graph) dfsCycle(v, parent int, visited []bool) bool {
    visited[v] = true
    for _, neighbor := range g.adj[v] {
        if !visited[neighbor] {
            if g.dfsCycle(neighbor, v, visited) {
                return true
            }
        } else if neighbor != parent {
            return true  // Back edge to non-parent = cycle
        }
    }
    return false
}
```

**Directed graph — DFS with recursion stack (3-color):**
```go
const (
    white = 0  // Unvisited
    grey  = 1  // In current DFS path (in recursion stack)
    black = 2  // Fully processed
)

func (g *Graph) HasDirectedCycle() bool {
    color := make([]int, g.vertices)
    for v := 0; v < g.vertices; v++ {
        if color[v] == white {
            if g.dfsCycleDirected(v, color) {
                return true
            }
        }
    }
    return false
}

func (g *Graph) dfsCycleDirected(v int, color []int) bool {
    color[v] = grey  // Mark as "in current path"

    for _, neighbor := range g.adj[v] {
        if color[neighbor] == grey {
            return true  // Back edge to grey node = cycle
        }
        if color[neighbor] == white && g.dfsCycleDirected(neighbor, color) {
            return true
        }
    }

    color[v] = black  // Done processing this vertex
    return false
}
```

### Quick Check
> 1. In directed cycle detection, what does "grey" mean?
> 2. Why is parent-tracking needed for undirected cycle detection?
> 3. What is a "back edge"?

---

## 5. Topological Sort

Topological sort orders vertices of a **directed acyclic graph (DAG)** so every edge points forward — if A→B, then A comes before B. Used for dependency resolution, build systems, course prerequisites.

**Kahn's Algorithm (BFS-based):**
```go
// TopologicalSortKahn returns vertices in topological order using BFS.
// Returns nil if the graph has a cycle.
func (g *Graph) TopologicalSortKahn() []int {
    inDegree := make([]int, g.vertices)
    for v := 0; v < g.vertices; v++ {
        for _, neighbor := range g.adj[v] {
            inDegree[neighbor]++
        }
    }

    // Start with all vertices that have no dependencies (in-degree 0):
    var queue []int
    for v, deg := range inDegree {
        if deg == 0 {
            queue = append(queue, v)
        }
    }

    var result []int
    for len(queue) > 0 {
        v := queue[0]
        queue = queue[1:]
        result = append(result, v)

        for _, neighbor := range g.adj[v] {
            inDegree[neighbor]--
            if inDegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }

    if len(result) != g.vertices {
        return nil  // Cycle detected (some vertices never reached in-degree 0)
    }
    return result
}
```

**DFS-based topological sort:**
```go
func (g *Graph) TopologicalSortDFS() []int {
    visited := make([]bool, g.vertices)
    var stack []int

    var dfs func(v int)
    dfs = func(v int) {
        visited[v] = true
        for _, neighbor := range g.adj[v] {
            if !visited[neighbor] {
                dfs(neighbor)
            }
        }
        stack = append(stack, v)  // Push AFTER processing children
    }

    for v := 0; v < g.vertices; v++ {
        if !visited[v] {
            dfs(v)
        }
    }

    // Reverse the stack:
    result := make([]int, len(stack))
    for i, v := range stack {
        result[len(stack)-1-i] = v
    }
    return result
}
```

**Real-world: package dependency resolution:**
```go
func ResolveDependencies(packages map[string][]string) ([]string, error) {
    // Build graph:
    inDegree := make(map[string]int)
    for pkg := range packages {
        if _, ok := inDegree[pkg]; !ok {
            inDegree[pkg] = 0
        }
        for _, dep := range packages[pkg] {
            if _, ok := inDegree[dep]; !ok {
                inDegree[dep] = 0 // Ensure dep exists in map
            }
            inDegree[pkg]++ // pkg depends on dep
        }
    }
    // ... rest of Kahn's algorithm with string keys ...
    return nil, nil
}
```

---

## 6. Shortest Path — BFS and Dijkstra

**BFS shortest path** works for unweighted graphs (Section 2 above).

**Dijkstra's algorithm** for weighted graphs — O((V+E) log V) with a priority queue:

```go
import "container/heap"

func (g *WeightedGraph) Dijkstra(start int) ([]int, []int) {
    dist := make([]int, g.vertices)
    parent := make([]int, g.vertices)
    for i := range dist {
        dist[i] = 1<<31 - 1  // MaxInt
        parent[i] = -1
    }
    dist[start] = 0

    // Min-heap of (vertex, tentative distance) — PQItemHeap defined below:
    h := &PQItemHeap{{start, 0}}
    heap.Init(h)

    for h.Len() > 0 {
        item := heap.Pop(h).(PQItem)
        v, d := item.v, item.d

        if d > dist[v] {
            continue  // Stale entry — skip
        }

        for _, edge := range g.adj[v] {
            newDist := dist[v] + edge.Weight
            if newDist < dist[edge.To] {
                dist[edge.To] = newDist
                parent[edge.To] = v
                heap.Push(h, PQItem{edge.To, newDist})
            }
        }
    }
    return dist, parent
}

// PQItemHeap implements heap.Interface for Dijkstra:
type PQItem struct{ v, d int }
type PQItemHeap []PQItem
func (h PQItemHeap) Len() int            { return len(h) }
func (h PQItemHeap) Less(i, j int) bool  { return h[i].d < h[j].d }
func (h PQItemHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *PQItemHeap) Push(x any)         { *h = append(*h, x.(PQItem)) }
func (h *PQItemHeap) Pop() any {
    old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x
}
```

---

## 7. Union-Find (Disjoint Set)

Union-Find efficiently answers: "are vertices X and Y in the same connected component?"

```go
type UnionFind struct {
    parent []int
    rank   []int
    count  int  // Number of disjoint sets
}

func NewUnionFind(n int) *UnionFind {
    parent := make([]int, n)
    rank := make([]int, n)
    for i := range parent { parent[i] = i }
    return &UnionFind{parent: parent, rank: rank, count: n}
}

// Find returns the root of the set containing x (with path compression):
func (uf *UnionFind) Find(x int) int {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x])  // Path compression
    }
    return uf.parent[x]
}

// Union merges the sets containing x and y (union by rank):
func (uf *UnionFind) Union(x, y int) bool {
    px, py := uf.Find(x), uf.Find(y)
    if px == py {
        return false  // Already in same set
    }
    if uf.rank[px] < uf.rank[py] {
        px, py = py, px
    }
    uf.parent[py] = px
    if uf.rank[px] == uf.rank[py] {
        uf.rank[px]++
    }
    uf.count--
    return true
}

func (uf *UnionFind) Connected(x, y int) bool {
    return uf.Find(x) == uf.Find(y)
}

func (uf *UnionFind) Count() int { return uf.count }
```

**Number of islands (classic):**
```go
func NumIslands(grid [][]byte) int {
    rows, cols := len(grid), len(grid[0])
    uf := NewUnionFind(rows * cols)
    islands := 0

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if grid[r][c] == '1' {
                islands++
                idx := r*cols + c
                // Connect to adjacent land cells:
                if r > 0 && grid[r-1][c] == '1' {
                    if uf.Union(idx, (r-1)*cols+c) { islands-- }
                }
                if c > 0 && grid[r][c-1] == '1' {
                    if uf.Union(idx, r*cols+(c-1)) { islands-- }
                }
            }
        }
    }
    return islands
}
```

**Time complexity:** With path compression + union by rank, each operation is nearly O(1) — specifically O(α(n)) where α is the inverse Ackermann function (< 5 for any practical n).

---

## Summary

- **Graph**: vertices + edges; directed/undirected, weighted/unweighted
- **Adjacency list**: O(V+E) space; best for sparse graphs (most real graphs)
- **Adjacency matrix**: O(V²) space; O(1) edge lookup; best for dense graphs
- **BFS**: queue-based; finds shortest path (unweighted); O(V+E)
- **DFS**: stack/recursive; finds paths, cycles, components; O(V+E)
- **Cycle detection**: undirected = parent tracking; directed = 3-color (white/grey/black)
- **Topological sort**: DAG only; Kahn's (BFS, detects cycle) or DFS-postorder
- **Dijkstra**: weighted shortest path; priority queue; O((V+E) log V)
- **Union-Find**: disjoint sets; Find+Union nearly O(1) with path compression + union by rank

---

## Exercises

### Easy
1. Implement `IsConnected(g *Graph) bool` — returns true if the graph has exactly one connected component (all vertices reachable from vertex 0).
2. Write `BipartiteCheck(g *Graph) bool` — a graph is bipartite if you can color all vertices with 2 colors such that no two adjacent vertices have the same color. Use BFS with alternating colors.
3. Write `CloneGraph(node *GraphNode) *GraphNode` — deep copy a graph where each node has a `Val` and `Neighbors []*GraphNode`. Use BFS and a map to track already-cloned nodes.

### Medium
4. Word ladder: Given a start word, end word, and word list, find the shortest transformation sequence where each step changes exactly one letter. E.g., `"hit" → "hot" → "dot" → "dog" → "cog"`. Build the graph implicitly (don't precompute all edges) using BFS. Return the number of words in the shortest sequence, or 0 if no path exists.
5. Course schedule: Given N courses (0 to N-1) and prerequisites `[[a,b]]` meaning "must take b before a", determine if you can finish all courses. Return the order to take them. Use topological sort (Kahn's). Test with: valid schedule, impossible schedule (cycle), and a large graph with 2000 courses.
6. Network delay time: Given a network of N nodes, N edges with weights, and a source node K, find the time for all nodes to receive a signal sent from K. Use Dijkstra. Return -1 if some node is unreachable.

### Hard
7. Minimum spanning tree (Kruskal's): Implement Kruskal's algorithm to find the MST of a weighted undirected graph. Sort edges by weight, use Union-Find to add edges that don't create a cycle. Return the total weight and the list of MST edges. Test with: a complete 5-node graph, disconnected graph (return -1 for impossible), graph where MST is unique vs multiple valid MSTs.
8. All-pairs shortest paths (Floyd-Warshall): Implement Floyd-Warshall for a weighted directed graph. `AllPairsShortestPaths(dist [][]int) [][]int` — `dist[i][j]` is the edge weight (MaxInt = no edge). After the algorithm, `dist[i][j]` should be the shortest path from i to j. Also detect negative cycles (shortest path = -∞). Test with a 5-node graph and verify all 25 pairs.
