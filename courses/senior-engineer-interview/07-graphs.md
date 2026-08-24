# Chapter 07: Graphs — BFS, DFS, Topological Sort & Union-Find

Graphs are the most versatile data structure in computer science. Networks, maps, dependencies, social connections — they are all graphs. Senior interviews frequently include graph problems because they test both algorithmic thinking and the ability to model real problems as graphs.

## Table of Contents

1. [Graph Representation in Go](#1-graph-representation-in-go)
2. [BFS on Graphs](#2-bfs-on-graphs)
3. [DFS on Graphs](#3-dfs-on-graphs)
4. [Topological Sort](#4-topological-sort)
5. [Union-Find (Disjoint Set)](#5-union-find-disjoint-set)
6. [Classic Graph Problems](#6-classic-graph-problems)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Graph Representation in Go

Two main representations: adjacency list (most common) and adjacency matrix.

```go
// ADJACENCY LIST — O(V+E) space, preferred for sparse graphs
// graph[v] = list of neighbors of v
type Graph struct {
    adj map[int][]int
}

func NewGraph() *Graph {
    return &Graph{adj: make(map[int][]int)}
}

func (g *Graph) AddEdge(u, v int) {
    g.adj[u] = append(g.adj[u], v)
    g.adj[v] = append(g.adj[v], u) // undirected: add both directions
}

// For interview problems, a simple map is often enough:
adj := map[int][]int{
    0: {1, 2},
    1: {0, 3},
    2: {0, 4},
    3: {1},
    4: {2},
}

// ADJACENCY MATRIX — O(V²) space, preferred for dense graphs or when checking
// edge existence frequently
matrix := make([][]int, n)
for i := range matrix { matrix[i] = make([]int, n) }
matrix[u][v] = 1 // add edge from u to v
```

### Building a Graph from Edge List

```go
// Common in interview problems: given n nodes and list of edges
func buildGraph(n int, edges [][]int) map[int][]int {
    adj := make(map[int][]int)
    for _, e := range edges {
        u, v := e[0], e[1]
        adj[u] = append(adj[u], v)
        adj[v] = append(adj[v], u)
    }
    return adj
}
```

---

## 2. BFS on Graphs

BFS explores level by level. It finds the **shortest path** in unweighted graphs. Always use a visited set to avoid revisiting nodes.

```go
func bfs(start int, adj map[int][]int) []int {
    visited := make(map[int]bool)
    queue := []int{start}
    visited[start] = true
    order := []int{}

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        order = append(order, node)

        for _, neighbor := range adj[node] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
    return order
}
// Time: O(V+E) — visit each vertex and each edge once
// Space: O(V) for visited set and queue
```

### Shortest Path in Unweighted Graph

```go
func shortestPath(start, end int, adj map[int][]int) int {
    if start == end { return 0 }
    
    visited := make(map[int]bool)
    queue := []int{start}
    visited[start] = true
    distance := 0

    for len(queue) > 0 {
        distance++
        levelSize := len(queue)
        for i := 0; i < levelSize; i++ {
            node := queue[0]
            queue = queue[1:]
            for _, neighbor := range adj[node] {
                if neighbor == end { return distance }
                if !visited[neighbor] {
                    visited[neighbor] = true
                    queue = append(queue, neighbor)
                }
            }
        }
    }
    return -1 // no path found
}
```

### Problem: Number of Islands (Grid BFS)

Grids are implicit graphs. Each cell is a node; edges connect adjacent cells.

```go
func numIslands(grid [][]byte) int {
    if len(grid) == 0 { return 0 }
    rows, cols := len(grid), len(grid[0])
    count := 0

    var bfs func(r, c int)
    bfs = func(r, c int) {
        if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] != '1' {
            return
        }
        grid[r][c] = '0' // mark visited by changing to water
        bfs(r+1, c); bfs(r-1, c); bfs(r, c+1); bfs(r, c-1)
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if grid[r][c] == '1' {
                count++
                bfs(r, c) // flood-fill the island
            }
        }
    }
    return count
}
// Time: O(rows * cols), Space: O(rows * cols) for DFS stack
// Note: modifying the input — if not allowed, use a separate visited array
```

### Problem: 0/1 Matrix — Shortest Distance to 0

```go
// Multi-source BFS: start from all 0s simultaneously
func updateMatrix(mat [][]int) [][]int {
    rows, cols := len(mat), len(mat[0])
    dist := make([][]int, rows)
    for i := range dist { dist[i] = make([]int, cols) }
    
    queue := [][]int{}
    dirs := [][]int{{0,1},{0,-1},{1,0},{-1,0}}

    // Initialize: all 0-cells have distance 0, all 1-cells have distance infinity
    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if mat[r][c] == 0 {
                queue = append(queue, []int{r, c})
                dist[r][c] = 0
            } else {
                dist[r][c] = 1<<31 - 1 // infinity
            }
        }
    }

    // BFS from all 0s simultaneously
    for len(queue) > 0 {
        cell := queue[0]; queue = queue[1:]
        r, c := cell[0], cell[1]
        for _, d := range dirs {
            nr, nc := r+d[0], c+d[1]
            if nr >= 0 && nr < rows && nc >= 0 && nc < cols {
                if dist[r][c]+1 < dist[nr][nc] {
                    dist[nr][nc] = dist[r][c] + 1
                    queue = append(queue, []int{nr, nc})
                }
            }
        }
    }
    return dist
}
```

---

## 3. DFS on Graphs

DFS explores as deep as possible before backtracking. Use it for connectivity, cycle detection, and component finding.

```go
func dfs(node int, visited map[int]bool, adj map[int][]int) {
    visited[node] = true
    for _, neighbor := range adj[node] {
        if !visited[neighbor] {
            dfs(neighbor, visited, adj)
        }
    }
}

// Count connected components
func countComponents(n int, edges [][]int) int {
    adj := buildGraph(n, edges)
    visited := make(map[int]bool)
    count := 0

    for i := 0; i < n; i++ {
        if !visited[i] {
            dfs(i, visited, adj)
            count++
        }
    }
    return count
}
```

### Cycle Detection in Directed Graph

```go
// Three states: 0=unvisited, 1=in current DFS path, 2=fully processed
// Cycle exists if we reach a node that is in the current path (state 1)
func hasCycleDirected(n int, adj map[int][]int) bool {
    state := make([]int, n) // 0=unvisited, 1=visiting, 2=visited

    var dfs func(node int) bool
    dfs = func(node int) bool {
        state[node] = 1 // mark as being visited (in current path)
        for _, neighbor := range adj[node] {
            if state[neighbor] == 1 { return true }  // back edge = cycle
            if state[neighbor] == 0 && dfs(neighbor) { return true }
        }
        state[node] = 2 // fully processed
        return false
    }

    for i := 0; i < n; i++ {
        if state[i] == 0 && dfs(i) { return true }
    }
    return false
}
```

---

## 4. Topological Sort

Topological sort orders nodes in a directed acyclic graph (DAG) such that for every edge (u → v), u comes before v. Classic use: dependency ordering, course schedules.

### Kahn's Algorithm (BFS-based)

```go
// Uses in-degree: count of edges pointing TO each node.
// Process nodes with in-degree 0 (no dependencies) first.
func topologicalSort(n int, adj map[int][]int) ([]int, bool) {
    indegree := make([]int, n)
    for u := range adj {
        for _, v := range adj[u] {
            indegree[v]++
        }
    }

    // Start with all nodes that have no dependencies
    queue := []int{}
    for i := 0; i < n; i++ {
        if indegree[i] == 0 {
            queue = append(queue, i)
        }
    }

    order := []int{}
    for len(queue) > 0 {
        node := queue[0]; queue = queue[1:]
        order = append(order, node)
        for _, neighbor := range adj[node] {
            indegree[neighbor]--
            if indegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }

    // If we processed all nodes, no cycle. Otherwise, cycle exists.
    return order, len(order) == n
}
```

### Problem: Course Schedule

```go
// Can you finish all courses given prerequisites?
// This is just cycle detection in a directed graph.
func canFinish(numCourses int, prerequisites [][]int) bool {
    adj := make(map[int][]int)
    for _, p := range prerequisites {
        adj[p[1]] = append(adj[p[1]], p[0]) // p[1] must come before p[0]
    }

    _, noCircle := topologicalSort(numCourses, adj)
    return noCircle
}

// Return the order to take courses
func findOrder(numCourses int, prerequisites [][]int) []int {
    adj := make(map[int][]int)
    for _, p := range prerequisites {
        adj[p[1]] = append(adj[p[1]], p[0])
    }

    order, ok := topologicalSort(numCourses, adj)
    if !ok { return nil }
    return order
}
```

---

## 5. Union-Find (Disjoint Set)

Union-Find efficiently answers "are these two nodes in the same component?" and merges components. Use it for connected component problems where you process edges one at a time.

```go
type UnionFind struct {
    parent []int
    rank   []int
    count  int // number of connected components
}

func NewUnionFind(n int) *UnionFind {
    uf := &UnionFind{
        parent: make([]int, n),
        rank:   make([]int, n),
        count:  n,
    }
    for i := range uf.parent { uf.parent[i] = i } // each node is its own parent
    return uf
}

// Find with path compression: flatten the tree for future queries
func (uf *UnionFind) Find(x int) int {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x]) // path compression
    }
    return uf.parent[x]
}

// Union by rank: attach smaller tree under larger tree
func (uf *UnionFind) Union(x, y int) bool {
    px, py := uf.Find(x), uf.Find(y)
    if px == py { return false } // already in the same component

    if uf.rank[px] < uf.rank[py] { px, py = py, px }
    uf.parent[py] = px
    if uf.rank[px] == uf.rank[py] { uf.rank[px]++ }
    uf.count--
    return true
}

func (uf *UnionFind) Connected(x, y int) bool {
    return uf.Find(x) == uf.Find(y)
}
// With path compression + union by rank: nearly O(1) per operation (amortized O(α(n)))
```

### Problem: Number of Connected Components

```go
func countComponentsUF(n int, edges [][]int) int {
    uf := NewUnionFind(n)
    for _, e := range edges {
        uf.Union(e[0], e[1])
    }
    return uf.count
}
```

### Problem: Detect Cycle in Undirected Graph

```go
// If two nodes are already in the same component when we process an edge,
// adding that edge creates a cycle.
func hasCycleUndirected(n int, edges [][]int) bool {
    uf := NewUnionFind(n)
    for _, e := range edges {
        if !uf.Union(e[0], e[1]) { // Union returns false if already connected
            return true // cycle detected
        }
    }
    return false
}
```

---

## 6. Classic Graph Problems

### Word Ladder (BFS Shortest Path on Implicit Graph)

```go
// Transform "hit" to "cog" changing one letter at a time, each word in wordList.
func ladderLength(beginWord string, endWord string, wordList []string) int {
    wordSet := make(map[string]bool)
    for _, w := range wordList { wordSet[w] = true }
    if !wordSet[endWord] { return 0 }

    queue := []string{beginWord}
    visited := map[string]bool{beginWord: true}
    steps := 1

    for len(queue) > 0 {
        steps++
        for size := len(queue); size > 0; size-- {
            word := queue[0]; queue = queue[1:]
            // Try changing each character
            for i := 0; i < len(word); i++ {
                for c := byte('a'); c <= byte('z'); c++ {
                    if c == word[i] { continue }
                    newWord := word[:i] + string(c) + word[i+1:]
                    if newWord == endWord { return steps }
                    if wordSet[newWord] && !visited[newWord] {
                        visited[newWord] = true
                        queue = append(queue, newWord)
                    }
                }
            }
        }
    }
    return 0
}
// Time: O(M² × N) where M=word length, N=word list size
```

### Clone Graph

```go
type Node struct {
    Val       int
    Neighbors []*Node
}

func cloneGraph(node *Node) *Node {
    if node == nil { return nil }
    cloned := make(map[*Node]*Node)

    var clone func(*Node) *Node
    clone = func(n *Node) *Node {
        if n == nil { return nil }
        if c, ok := cloned[n]; ok { return c }

        c := &Node{Val: n.Val}
        cloned[n] = c // register before recursing (handles cycles)
        for _, neighbor := range n.Neighbors {
            c.Neighbors = append(c.Neighbors, clone(neighbor))
        }
        return c
    }
    return clone(node)
}
```

### Pacific Atlantic Water Flow

```go
// Multi-source BFS from both oceans. A cell can flow to both oceans if
// it is reachable from the Pacific BFS AND the Atlantic BFS.
func pacificAtlantic(heights [][]int) [][]int {
    rows, cols := len(heights), len(heights[0])
    dirs := [][]int{{0,1},{0,-1},{1,0},{-1,0}}

    bfsFromEdge := func(seeds [][]int) [][]bool {
        reachable := make([][]bool, rows)
        for i := range reachable { reachable[i] = make([]bool, cols) }
        queue := seeds
        for _, s := range seeds { reachable[s[0]][s[1]] = true }

        for len(queue) > 0 {
            cell := queue[0]; queue = queue[1:]
            r, c := cell[0], cell[1]
            for _, d := range dirs {
                nr, nc := r+d[0], c+d[1]
                if nr < 0 || nr >= rows || nc < 0 || nc >= cols { continue }
                if reachable[nr][nc] { continue }
                if heights[nr][nc] < heights[r][c] { continue } // water flows downhill
                reachable[nr][nc] = true
                queue = append(queue, []int{nr, nc})
            }
        }
        return reachable
    }

    // Pacific: top row + left column
    pacificSeeds := [][]int{}
    for r := 0; r < rows; r++ { pacificSeeds = append(pacificSeeds, []int{r, 0}) }
    for c := 1; c < cols; c++ { pacificSeeds = append(pacificSeeds, []int{0, c}) }

    // Atlantic: bottom row + right column
    atlanticSeeds := [][]int{}
    for r := 0; r < rows; r++ { atlanticSeeds = append(atlanticSeeds, []int{r, cols-1}) }
    for c := 0; c < cols-1; c++ { atlanticSeeds = append(atlanticSeeds, []int{rows-1, c}) }

    pacific := bfsFromEdge(pacificSeeds)
    atlantic := bfsFromEdge(atlanticSeeds)

    result := [][]int{}
    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if pacific[r][c] && atlantic[r][c] {
                result = append(result, []int{r, c})
            }
        }
    }
    return result
}
```

---

## Summary

- **Adjacency list:** default graph representation, O(V+E) space.
- **BFS:** shortest path in unweighted graphs. Queue + visited set. O(V+E).
- **DFS:** connectivity, cycle detection, component counting. Stack/recursion + visited. O(V+E).
- **Topological sort (Kahn's):** process zero in-degree nodes first. Detects cycles. O(V+E).
- **Union-Find:** dynamic connectivity. Nearly O(1) per operation with path compression + union by rank.
- Grid problems are graph problems: treat each cell as a node, adjacent cells as edges.
- Multi-source BFS: start from multiple sources simultaneously for problems like "shortest distance to nearest X."

---

## Exercises

### Easy
1. Find if a path exists between source and destination in an undirected graph.
2. Count the number of islands using DFS (modify the grid to mark visited).
3. Determine if a graph is a valid tree (connected + no cycles with exactly n-1 edges).

### Medium
4. Find all nodes in a connected component of a given starting node.
5. Determine if a graph is bipartite (can be colored with 2 colors such that no two adjacent nodes have the same color). Use BFS/DFS.
6. Given a list of accounts where each entry is [name, email1, email2, ...], merge accounts that share an email address. Return merged accounts. (Union-Find)

### Hard
7. Implement Dijkstra's algorithm from scratch in Go (covered in detail in the next chapter, but try it here first).
8. Given a grid with walls (0) and open space (1), find the shortest path from top-left to bottom-right where you can eliminate at most k walls.
9. Find all strongly connected components in a directed graph using Tarjan's or Kosaraju's algorithm.
