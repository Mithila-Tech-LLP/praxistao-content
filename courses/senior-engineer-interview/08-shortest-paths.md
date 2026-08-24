# Chapter 08: Shortest Paths & Advanced Graph Algorithms

Dijkstra's algorithm is the most commonly tested shortest-path algorithm in senior interviews. You must be able to implement it from scratch using a priority queue in Go. This chapter also covers Bellman-Ford for negative edges and Union-Find based MST algorithms.

## Table of Contents

1. [Dijkstra's Algorithm](#1-dijkstras-algorithm)
2. [Bellman-Ford Algorithm](#2-bellman-ford-algorithm)
3. [Minimum Spanning Tree — Kruskal's & Prim's](#3-minimum-spanning-tree)
4. [A* Search (Conceptual)](#4-a-search)
5. [Classic Problems](#5-classic-problems)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. Dijkstra's Algorithm

### The Core Idea

Dijkstra finds the shortest path from a source node to all other nodes in a **weighted graph with non-negative edges**. It works like BFS but with a priority queue — always process the unvisited node with the smallest known distance.

```
Algorithm:
1. Initialize dist[source] = 0, all others = infinity
2. Push source into min-heap with priority 0
3. While heap is non-empty:
   a. Pop node with minimum distance
   b. If already visited, skip (we found a shorter path earlier)
   c. Mark as visited
   d. For each neighbor:
      - newDist = dist[current] + edge weight
      - If newDist < dist[neighbor], update and push to heap
4. Return dist array
```

### Implementation in Go

```go
import "container/heap"

// Priority queue item
type Item struct {
    node, dist int
}

// Min-heap by distance
type PQ []Item
func (h PQ) Len() int            { return len(h) }
func (h PQ) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h PQ) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *PQ) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *PQ) Pop() interface{} {
    old := *h; n := len(old)
    x := old[n-1]; *h = old[:n-1]
    return x
}

// Weighted graph: adj[u] = list of {neighbor, weight}
type Edge struct{ to, weight int }

func dijkstra(src int, adj map[int][]Edge, n int) []int {
    dist := make([]int, n)
    for i := range dist { dist[i] = 1<<31 - 1 } // infinity
    dist[src] = 0

    pq := &PQ{{src, 0}}
    heap.Init(pq)
    visited := make([]bool, n)

    for pq.Len() > 0 {
        item := heap.Pop(pq).(Item)
        u := item.node

        if visited[u] { continue } // already processed with shorter path
        visited[u] = true

        for _, e := range adj[u] {
            newDist := dist[u] + e.weight
            if newDist < dist[e.to] {
                dist[e.to] = newDist
                heap.Push(pq, Item{e.to, newDist})
            }
        }
    }
    return dist
}
// Time: O((V+E) log V) — each edge may cause a heap push
// Space: O(V+E)
```

### Problem: Network Delay Time

```go
// Find the time for all nodes to receive a signal from source k.
// This is Dijkstra from source k, return max of all distances.
func networkDelayTime(times [][]int, n int, k int) int {
    adj := make(map[int][]Edge)
    for _, t := range times {
        adj[t[0]] = append(adj[t[0]], Edge{t[1], t[2]})
    }

    dist := dijkstra(k, adj, n+1) // nodes are 1-indexed

    maxDist := 0
    for i := 1; i <= n; i++ {
        if dist[i] == 1<<31-1 { return -1 } // node unreachable
        if dist[i] > maxDist { maxDist = dist[i] }
    }
    return maxDist
}
```

### Problem: Cheapest Flights Within K Stops

This is a variant: at most K stops (K+1 edges). Cannot use standard Dijkstra — use Bellman-Ford with K+1 iterations instead.

```go
func findCheapestPrice(n int, flights [][]int, src, dst, k int) int {
    // dist[i] = minimum cost to reach node i with at most 'stops' stops used so far
    dist := make([]int, n)
    for i := range dist { dist[i] = 1<<31 - 1 }
    dist[src] = 0

    // Run k+1 iterations of Bellman-Ford (k stops = k+1 edges)
    for i := 0; i <= k; i++ {
        // Important: copy dist to avoid using edges from the current iteration
        temp := make([]int, n)
        copy(temp, dist)

        for _, f := range flights {
            u, v, w := f[0], f[1], f[2]
            if dist[u] != 1<<31-1 && dist[u]+w < temp[v] {
                temp[v] = dist[u] + w
            }
        }
        dist = temp
    }

    if dist[dst] == 1<<31-1 { return -1 }
    return dist[dst]
}
// Time: O(k * E), Space: O(n)
```

---

## 2. Bellman-Ford Algorithm

Bellman-Ford finds shortest paths from a source, handles **negative edge weights**, and detects **negative cycles**. Slower than Dijkstra: O(VE).

```go
func bellmanFord(src int, edges [][]int, n int) ([]int, bool) {
    dist := make([]int, n)
    for i := range dist { dist[i] = 1<<31 - 1 }
    dist[src] = 0

    // Relax all edges V-1 times (shortest path has at most V-1 edges)
    for i := 0; i < n-1; i++ {
        for _, e := range edges {
            u, v, w := e[0], e[1], e[2]
            if dist[u] != 1<<31-1 && dist[u]+w < dist[v] {
                dist[v] = dist[u] + w
            }
        }
    }

    // Check for negative cycle: if we can still relax on the nth iteration,
    // a negative cycle exists
    for _, e := range edges {
        u, v, w := e[0], e[1], e[2]
        if dist[u] != 1<<31-1 && dist[u]+w < dist[v] {
            return nil, true // negative cycle detected
        }
    }
    return dist, false
}
// Time: O(VE), Space: O(V)
```

### When to Use Bellman-Ford vs Dijkstra

| Criteria | Dijkstra | Bellman-Ford |
|---|---|---|
| Edge weights | Non-negative only | Any (including negative) |
| Negative cycles | Cannot handle | Detects them |
| Time complexity | O((V+E) log V) | O(VE) |
| Space | O(V+E) | O(V) |
| Use when | Most common cases | Negative edges, cycle detection |

---

## 3. Minimum Spanning Tree

An MST connects all vertices with minimum total edge weight and no cycles.

### Kruskal's Algorithm (Union-Find)

```go
// Sort edges by weight. Add edge if it doesn't create a cycle.
func kruskal(n int, edges [][]int) int {
    // Sort edges by weight
    sort.Slice(edges, func(i, j int) bool {
        return edges[i][2] < edges[j][2]
    })

    uf := NewUnionFind(n)
    totalWeight := 0
    edgesUsed := 0

    for _, e := range edges {
        u, v, w := e[0], e[1], e[2]
        if uf.Union(u, v) { // only add if they are in different components
            totalWeight += w
            edgesUsed++
            if edgesUsed == n-1 { break } // MST has n-1 edges
        }
    }
    return totalWeight
}
// Time: O(E log E) for sorting, nearly O(E) for union-find operations
```

### Problem: Min Cost to Connect All Points

```go
// Points are nodes. Cost to connect = Manhattan distance.
// Find MST of this complete graph.
// Prim's algorithm is more efficient here — don't need to enumerate all edges.
func minCostConnectPoints(points [][]int) int {
    n := len(points)
    inMST := make([]bool, n)
    minCost := make([]int, n)
    for i := range minCost { minCost[i] = 1<<31 - 1 }
    minCost[0] = 0

    total := 0
    for i := 0; i < n; i++ {
        // Pick the node not in MST with the minimum edge cost
        u := -1
        for v := 0; v < n; v++ {
            if !inMST[v] && (u == -1 || minCost[v] < minCost[u]) {
                u = v
            }
        }
        inMST[u] = true
        total += minCost[u]

        // Update costs for neighbors
        for v := 0; v < n; v++ {
            if !inMST[v] {
                dist := abs(points[u][0]-points[v][0]) + abs(points[u][1]-points[v][1])
                if dist < minCost[v] {
                    minCost[v] = dist
                }
            }
        }
    }
    return total
}
// Time: O(n²) — Prim's with adjacency matrix; use priority queue for O(E log V)
```

---

## 4. A* Search

A* is Dijkstra + a heuristic. It finds shortest paths faster when you have domain knowledge about the goal direction.

The heuristic h(n) estimates distance from node n to the goal. Priority = actual cost g(n) + estimated cost h(n).

```
f(n) = g(n) + h(n)
g(n) = actual cost from source to n
h(n) = heuristic estimate from n to goal
```

For grids: Manhattan distance is a common admissible heuristic.

```go
// Conceptual A* implementation for a grid
func aStar(grid [][]int, start, end [2]int) int {
    rows, cols := len(grid), len(grid[0])
    dirs := [][2]int{{0,1},{0,-1},{1,0},{-1,0}}

    heuristic := func(r, c int) int {
        return abs(r-end[0]) + abs(c-end[1]) // Manhattan distance
    }

    type State struct{ f, g, r, c int }
    pq := []State{{heuristic(start[0], start[1]), 0, start[0], start[1]}}
    dist := make([][]int, rows)
    for i := range dist { dist[i] = make([]int, cols); for j := range dist[i] { dist[i][j] = 1<<31-1 } }
    dist[start[0]][start[1]] = 0

    for len(pq) > 0 {
        // Pop minimum f-score (simplified: use sort for interview, heap in production)
        sort.Slice(pq, func(i, j int) bool { return pq[i].f < pq[j].f })
        state := pq[0]; pq = pq[1:]
        r, c, g := state.r, state.c, state.g

        if r == end[0] && c == end[1] { return g }
        if g > dist[r][c] { continue } // stale entry

        for _, d := range dirs {
            nr, nc := r+d[0], c+d[1]
            if nr < 0 || nr >= rows || nc < 0 || nc >= cols || grid[nr][nc] == 1 { continue }
            ng := g + 1
            if ng < dist[nr][nc] {
                dist[nr][nc] = ng
                pq = append(pq, State{ng + heuristic(nr, nc), ng, nr, nc})
            }
        }
    }
    return -1
}
```

---

## 5. Classic Problems

### Path with Maximum Probability

```go
// Modified Dijkstra: maximize probability instead of minimize distance.
// Use a max-heap (negate probability for Go's min-heap).
func maxProbability(n int, edges [][]int, succProb []float64, start, end int) float64 {
    adj := make([][]struct{ to int; prob float64 }, n)
    for i, e := range edges {
        adj[e[0]] = append(adj[e[0]], struct{ to int; prob float64 }{e[1], succProb[i]})
        adj[e[1]] = append(adj[e[1]], struct{ to int; prob float64 }{e[0], succProb[i]})
    }

    prob := make([]float64, n)
    prob[start] = 1.0

    // Use a max-heap (sort by negative probability)
    type Item struct{ node int; p float64 }
    pq := []Item{{start, 1.0}}

    for len(pq) > 0 {
        sort.Slice(pq, func(i, j int) bool { return pq[i].p > pq[j].p })
        curr := pq[0]; pq = pq[1:]

        if curr.node == end { return curr.p }
        if curr.p < prob[curr.node] { continue }

        for _, e := range adj[curr.node] {
            newProb := prob[curr.node] * e.prob
            if newProb > prob[e.to] {
                prob[e.to] = newProb
                pq = append(pq, Item{e.to, newProb})
            }
        }
    }
    return prob[end]
}
```

---

## Summary

- **Dijkstra:** shortest path, non-negative weights only. O((V+E) log V) with priority queue. Know the implementation cold.
- **Bellman-Ford:** handles negative weights, detects negative cycles. O(VE) — slow but correct.
- **K-stops constraint:** use Bellman-Ford with K+1 iterations, not Dijkstra.
- **Kruskal's MST:** sort edges + Union-Find. O(E log E).
- **Prim's MST:** grow MST from a seed using greedy min-cost expansion. O(V²) naive, O(E log V) with heap.
- **A*:** Dijkstra + heuristic. Faster for spatial problems when you know direction to goal.

---

## Exercises

### Easy
1. Find if a path exists with total weight ≤ budget from source to destination using Dijkstra.
2. Find the maximum edge in the minimum bottleneck path between two nodes.

### Medium
3. Implement Dijkstra's algorithm without using Go's `container/heap` package — use a sorted slice as the priority queue. Compare the complexity.
4. Find the city with the smallest number of neighbors within distance threshold (weighted graph). Use Dijkstra from each node.
5. Given a directed weighted graph, find if there is a negative cycle reachable from the source.

### Hard
6. Implement Floyd-Warshall algorithm for all-pairs shortest paths in O(V³). When is this better than running Dijkstra from every vertex?
7. Implement the "Swim in Rising Water" problem — find the minimum time to reach (n-1, n-1) from (0,0) where the water level rises each second. (Modified Dijkstra on grid)
8. Find the shortest path in a grid where you can eliminate at most k obstacles. (Dijkstra with state = (row, col, k_remaining))
