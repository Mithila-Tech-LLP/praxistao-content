# Chapter 54: Shortest Path — Bellman-Ford and Floyd-Warshall

Dijkstra's algorithm (Ch 53) is the fastest single-source shortest path algorithm for graphs with non-negative weights. But it fails with negative edge weights. Bellman-Ford handles negative weights and detects negative cycles. Floyd-Warshall finds shortest paths between all pairs of vertices.

## Table of Contents

1. [Bellman-Ford](#1-bellman-ford)
2. [Negative Cycle Detection](#2-negative-cycle-detection)
3. [Floyd-Warshall](#3-floyd-warshall)
4. [Algorithm Comparison](#4-algorithm-comparison)
5. [When to Use Each](#5-when-to-use-each)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. Bellman-Ford

Bellman-Ford finds the shortest path from a single source to all vertices, even with negative edge weights. It relaxes **all edges** n-1 times.

**Why n-1 times?** A shortest path in a graph with n vertices has at most n-1 edges (a simple path visits each vertex at most once). After k relaxation rounds, we have correct distances for all paths using at most k edges.

```go
package bellman

import "math"

type Edge struct {
    From, To, Weight int
}

type Graph struct {
    Vertices int
    Edges    []Edge
}

// ShortestPath returns distances from src to all vertices.
// Returns (distances, true) if no negative cycle is reachable from src.
// Returns (nil, false) if a negative cycle is detected.
func ShortestPath(g *Graph, src int) ([]int, bool) {
    dist := make([]int, g.Vertices)
    for i := range dist { dist[i] = math.MaxInt64 }
    dist[src] = 0

    // Relax all edges n-1 times
    for i := 0; i < g.Vertices-1; i++ {
        updated := false
        for _, e := range g.Edges {
            if dist[e.From] != math.MaxInt64 &&
                dist[e.From]+e.Weight < dist[e.To] {
                dist[e.To] = dist[e.From] + e.Weight
                updated = true
            }
        }
        if !updated { break } // early exit: no changes, done
    }

    // Detect negative cycles: if we can still relax, there's a negative cycle
    for _, e := range g.Edges {
        if dist[e.From] != math.MaxInt64 &&
            dist[e.From]+e.Weight < dist[e.To] {
            return nil, false // negative cycle detected
        }
    }

    return dist, true
}
```

### Trace example

```
Vertices: 0, 1, 2, 3
Edges: 0→1 (weight 4), 0→2 (weight 5), 1→2 (-3), 2→3 (2), 1→3 (6)
Source: 0

Initial: dist = [0, ∞, ∞, ∞]

Round 1 (relax all edges):
  0→1: dist[1] = 0+4 = 4         dist = [0, 4, ∞, ∞]
  0→2: dist[2] = 0+5 = 5         dist = [0, 4, 5, ∞]
  1→2: dist[2] = min(5, 4-3) = 1  dist = [0, 4, 1, ∞]
  2→3: dist[3] = 1+2 = 3         dist = [0, 4, 1, 3]
  1→3: dist[3] = min(3, 4+6) = 3  no change

Round 2: No changes (no relaxation succeeds) → done early

Final: dist = [0, 4, 1, 3]
```

### Path reconstruction

```go
type BellmanFordResult struct {
    Dist   []int
    Parent []int
}

func ShortestPathWithParent(g *Graph, src int) (*BellmanFordResult, bool) {
    dist := make([]int, g.Vertices)
    parent := make([]int, g.Vertices)
    for i := range dist {
        dist[i] = math.MaxInt64
        parent[i] = -1
    }
    dist[src] = 0

    for i := 0; i < g.Vertices-1; i++ {
        for _, e := range g.Edges {
            if dist[e.From] != math.MaxInt64 &&
                dist[e.From]+e.Weight < dist[e.To] {
                dist[e.To] = dist[e.From] + e.Weight
                parent[e.To] = e.From
            }
        }
    }

    for _, e := range g.Edges {
        if dist[e.From] != math.MaxInt64 &&
            dist[e.From]+e.Weight < dist[e.To] {
            return nil, false
        }
    }

    return &BellmanFordResult{Dist: dist, Parent: parent}, true
}

func (r *BellmanFordResult) Path(dst int) []int {
    var path []int
    for v := dst; v != -1; v = r.Parent[v] {
        path = append(path, v)
    }
    // Reverse
    for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
        path[i], path[j] = path[j], path[i]
    }
    return path
}
```

---

## 2. Negative Cycle Detection

Bellman-Ford doesn't just detect negative cycles — it can identify which vertices are affected by them.

```go
// MarkNegativeCycleAffected marks all vertices with distance = -∞
// (i.e., affected by a negative cycle reachable from src)
func MarkNegativeCycleAffected(g *Graph, src int) []int {
    dist := make([]int, g.Vertices)
    for i := range dist { dist[i] = math.MaxInt64 }
    dist[src] = 0

    for i := 0; i < g.Vertices-1; i++ {
        for _, e := range g.Edges {
            if dist[e.From] != math.MaxInt64 &&
                dist[e.From]+e.Weight < dist[e.To] {
                dist[e.To] = dist[e.From] + e.Weight
            }
        }
    }

    // Mark with -∞ using BFS/DFS from nodes in negative cycles
    inNegCycle := make([]bool, g.Vertices)
    for _, e := range g.Edges {
        if dist[e.From] != math.MaxInt64 &&
            dist[e.From]+e.Weight < dist[e.To] {
            inNegCycle[e.To] = true
        }
    }

    // Propagate: anything reachable from a negative cycle vertex also gets -∞
    // BFS over edge list
    changed := true
    for changed {
        changed = false
        for _, e := range g.Edges {
            if inNegCycle[e.From] && !inNegCycle[e.To] {
                inNegCycle[e.To] = true
                changed = true
            }
        }
    }

    const negInf = math.MinInt64
    for i, neg := range inNegCycle {
        if neg { dist[i] = negInf }
    }
    return dist
}
```

### Real-world use: currency arbitrage

A negative cycle in a graph of currency exchange rates means you can make infinite profit by cycling through currencies. Bellman-Ford detects this:

```go
// Convert to log-space to use additive distances
// log(1/rate) = -log(rate)
// A negative cycle in -log(rate) weights = arbitrage opportunity
func detectArbitrage(currencies []string, rates [][]float64) bool {
    n := len(currencies)
    edges := []Edge{}
    
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if i != j && rates[i][j] > 0 {
                // Weight = -log(rate): negative cycle = positive log-product = arbitrage
                w := int(-math.Log(rates[i][j]) * 1e9) // scale to int
                edges = append(edges, Edge{From: i, To: j, Weight: w})
            }
        }
    }
    
    g := &Graph{Vertices: n, Edges: edges}
    _, ok := ShortestPath(g, 0)
    return !ok // not-ok means negative cycle = arbitrage
}
```

---

## 3. Floyd-Warshall

Floyd-Warshall finds shortest paths between **all pairs** of vertices in O(n³) time and O(n²) space. It works with negative weights but not negative cycles.

**Key insight**: `dist[i][j]` via intermediate vertex `k` = `dist[i][k] + dist[k][j]`. Iterate over all possible intermediate vertices k, updating the distance matrix.

```go
package floyd

import "math"

const INF = math.MaxInt64 / 2 // avoid overflow when adding

// AllPairsShortestPath returns the distance matrix and parent matrix.
// dist[i][j] = shortest distance from i to j.
// parent[i][j] = the vertex before j on the path from i to j.
func AllPairsShortestPath(n int, edges [][3]int) (dist, parent [][]int, hasNegCycle bool) {
    // Initialize distance matrix
    dist = make([][]int, n)
    parent = make([][]int, n)
    for i := range dist {
        dist[i] = make([]int, n)
        parent[i] = make([]int, n)
        for j := range dist[i] {
            if i == j {
                dist[i][j] = 0
            } else {
                dist[i][j] = INF
            }
            parent[i][j] = -1
        }
    }

    // Initialize with direct edges
    for _, e := range edges {
        u, v, w := e[0], e[1], e[2]
        if w < dist[u][v] {
            dist[u][v] = w
            parent[u][v] = u
        }
    }

    // Relax via all intermediate vertices
    for k := 0; k < n; k++ {
        for i := 0; i < n; i++ {
            for j := 0; j < n; j++ {
                if dist[i][k] == INF || dist[k][j] == INF {
                    continue
                }
                if dist[i][k]+dist[k][j] < dist[i][j] {
                    dist[i][j] = dist[i][k] + dist[k][j]
                    parent[i][j] = parent[k][j]
                }
            }
        }
    }

    // Detect negative cycles: if any dist[i][i] < 0, negative cycle exists
    for i := 0; i < n; i++ {
        if dist[i][i] < 0 {
            hasNegCycle = true
            return
        }
    }

    return dist, parent, false
}

// ReconstructPath reconstructs the path from src to dst.
func ReconstructPath(parent [][]int, src, dst int) []int {
    if parent[src][dst] == -1 {
        return nil // no path
    }
    var path []int
    for v := dst; v != src; v = parent[src][v] {
        path = append(path, v)
    }
    path = append(path, src)
    for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
        path[i], path[j] = path[j], path[i]
    }
    return path
}
```

### Transitive closure

Floyd-Warshall also computes whether any path exists between vertices (transitive closure):

```go
// TransitiveClosure[i][j] = true if j is reachable from i
func TransitiveClosure(n int, edges [][2]int) [][]bool {
    reach := make([][]bool, n)
    for i := range reach {
        reach[i] = make([]bool, n)
        reach[i][i] = true
    }
    for _, e := range edges { reach[e[0]][e[1]] = true }
    
    for k := 0; k < n; k++ {
        for i := 0; i < n; i++ {
            for j := 0; j < n; j++ {
                if reach[i][k] && reach[k][j] {
                    reach[i][j] = true
                }
            }
        }
    }
    return reach
}
```

---

## 4. Algorithm Comparison

| Algorithm | Weights | Time | Space | Single/All |
|-----------|---------|------|-------|------------|
| BFS | Unweighted | O(V+E) | O(V) | Single source |
| Dijkstra | Non-negative | O((V+E) log V) | O(V) | Single source |
| Bellman-Ford | Any (detects neg cycle) | O(VE) | O(V) | Single source |
| Floyd-Warshall | Any (no neg cycle) | O(V³) | O(V²) | All pairs |
| Johnson's | Any (no neg cycle) | O(V² log V + VE) | O(V²) | All pairs |

**When Floyd-Warshall over Dijkstra × V:**
- Dense graphs: Dijkstra × V = O(V³ log V), Floyd is O(V³) — Floyd wins
- Sparse graphs: Dijkstra × V = O(V(V+E) log V), Floyd is O(V³) — Dijkstra wins

---

## 5. When to Use Each

```
Has negative edges?
  No  → Dijkstra (fastest for single source)
  Yes → Need single source?
          Yes → Bellman-Ford
          No  → Need all pairs?
                  Yes → Floyd-Warshall (or Johnson's for sparse)
                  No  → Bellman-Ford

Dense graph AND all pairs?
  → Floyd-Warshall

Sparse graph AND all pairs?
  → Run Dijkstra from each vertex (or Johnson's algorithm)
```

---

## Summary

- **Bellman-Ford**: O(VE), handles negative weights, detects negative cycles. Use when Dijkstra can't (negative weights).
- **Floyd-Warshall**: O(V³), O(V²) space, finds all-pairs shortest paths. Simple 3-loop DP.
- Both use **relaxation**: if `dist[u] + w(u,v) < dist[v]`, update `dist[v]`.
- Bellman-Ford relaxes all edges n-1 times. Floyd-Warshall tries each vertex as an intermediate.
- Negative cycle detection: Bellman-Ford — one more relaxation round still succeeds. Floyd-Warshall — any `dist[i][i] < 0`.

---

## Exercises

### Easy
1. Trace Bellman-Ford on: 4 vertices, edges: `0→1(1)`, `1→2(-2)`, `2→3(1)`, `0→3(10)`. Show dist array after each round.
2. Trace Floyd-Warshall on a 3-vertex graph with edges `0→1(3)`, `1→2(-2)`, `0→2(4)`. Show the distance matrix after each value of k.
3. Implement `shortestCycle(n int, edges [][3]int) int` using Floyd-Warshall: the shortest cycle in the graph is `min(dist[i][j] + w(j,i))` for all edges (j,i).

### Medium
4. Implement **SPFA (Shortest Path Faster Algorithm)**: Bellman-Ford with a queue — only add vertex `v` to the queue if `dist[v]` was just updated. This reduces average-case complexity to O(E) while keeping worst-case O(VE). Benchmark on random vs adversarial graphs.
5. Implement **Johnson's algorithm**: reweight the graph using Bellman-Ford so all edges are non-negative, then run Dijkstra from each vertex. This gives O(V² log V + VE) for all-pairs on sparse graphs.
6. **Cheapest flights within k stops**: given flights `(from, to, price)` and a budget of k stops, find the cheapest flight from src to dst. Use a modified Bellman-Ford where you only relax k+1 times.

### Hard
7. Implement **minimum mean cycle**: given a directed weighted graph, find the cycle with the minimum average edge weight. Use Floyd-Warshall or a dynamic programming approach based on Bellman-Ford to solve in O(V³) time.
8. Use Bellman-Ford to solve the **longest path problem in a DAG** by negating all edge weights. Implement a function that finds the longest path from a source vertex in a DAG, and verify that converting to negative weights makes Bellman-Ford solve it correctly.
