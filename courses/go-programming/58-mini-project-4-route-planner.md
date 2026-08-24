# Chapter 58: Mini Project 4 — Route Planner with Dijkstra

A command-line route planner that builds a weighted graph from a JSON city map, finds shortest paths using Dijkstra's algorithm, and outputs turn-by-turn directions. This project ties together graphs, heaps (priority queue), and Bellman-Ford for negative-edge detection.

## What You'll Build

```
routeplanner/
├── main.go
├── graph/
│   ├── graph.go     # Weighted directed graph
│   └── loader.go    # Load from JSON file
├── algo/
│   ├── dijkstra.go  # O((V+E) log V) shortest path
│   └── bellman.go   # Negative-weight detection
└── maps/
    └── cities.json  # Sample city map
```

---

## 1. City Map Format

```json
{
  "cities": ["A", "B", "C", "D", "E"],
  "roads": [
    {"from": "A", "to": "B", "km": 100, "name": "Highway 1"},
    {"from": "A", "to": "C", "km": 150, "name": "Route 66"},
    {"from": "B", "to": "C", "km":  40, "name": "Main St"},
    {"from": "B", "to": "D", "km":  80, "name": "Park Ave"},
    {"from": "C", "to": "D", "km":  60, "name": "Oak Rd"},
    {"from": "C", "to": "E", "km": 120, "name": "Coastal Hwy"},
    {"from": "D", "to": "E", "km":  50, "name": "Lake Dr"}
  ]
}
```

---

## 2. Graph

```go
// graph/graph.go
package graph

type Edge struct {
    To     int
    Weight int
    Name   string
}

type Graph struct {
    Cities []string           // city ID → name
    Adj    [][]Edge           // adjacency list
    Index  map[string]int     // city name → ID
}

func New() *Graph {
    return &Graph{Index: make(map[string]int)}
}

func (g *Graph) AddCity(name string) int {
    if id, ok := g.Index[name]; ok {
        return id
    }
    id := len(g.Cities)
    g.Cities = append(g.Cities, name)
    g.Adj = append(g.Adj, nil)
    g.Index[name] = id
    return id
}

func (g *Graph) AddRoad(from, to string, km int, roadName string) {
    f := g.AddCity(from)
    t := g.AddCity(to)
    g.Adj[f] = append(g.Adj[f], Edge{To: t, Weight: km, Name: roadName})
    // Undirected: add reverse too
    g.Adj[t] = append(g.Adj[t], Edge{To: f, Weight: km, Name: roadName})
}

func (g *Graph) NumVertices() int { return len(g.Cities) }
```

```go
// graph/loader.go
package graph

import (
    "encoding/json"
    "fmt"
    "os"
)

type mapFile struct {
    Cities []string `json:"cities"`
    Roads  []struct {
        From string `json:"from"`
        To   string `json:"to"`
        KM   int    `json:"km"`
        Name string `json:"name"`
    } `json:"roads"`
}

func LoadFromFile(path string) (*Graph, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read map file: %w", err)
    }
    var mf mapFile
    if err := json.Unmarshal(data, &mf); err != nil {
        return nil, fmt.Errorf("parse map file: %w", err)
    }
    g := New()
    for _, city := range mf.Cities {
        g.AddCity(city)
    }
    for _, road := range mf.Roads {
        if road.KM < 0 {
            return nil, fmt.Errorf("road %s→%s has negative distance %d (use Bellman-Ford for this)",
                road.From, road.To, road.KM)
        }
        g.AddRoad(road.From, road.To, road.KM, road.Name)
    }
    return g, nil
}
```

---

## 3. Dijkstra

```go
// algo/dijkstra.go
package algo

import (
    "container/heap"
    "math"

    "routeplanner/graph"
)

type pqItem struct {
    vertex int
    dist   int
    index  int
}

type priorityQueue []*pqItem

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}
func (pq *priorityQueue) Push(x any) {
    item := x.(*pqItem)
    item.index = len(*pq)
    *pq = append(*pq, item)
}
func (pq *priorityQueue) Pop() any {
    old := *pq
    n := len(old)
    item := old[n-1]
    old[n-1] = nil
    item.index = -1
    *pq = old[:n-1]
    return item
}

type Result struct {
    Dist   []int
    Parent []int // parent[v] = the vertex before v on the shortest path from src
    Road   []string // road[v] = name of road used to reach v
}

func Dijkstra(g *graph.Graph, src int) *Result {
    n := g.NumVertices()
    dist := make([]int, n)
    parent := make([]int, n)
    road := make([]string, n)
    for i := range dist {
        dist[i] = math.MaxInt64
        parent[i] = -1
    }
    dist[src] = 0

    pq := &priorityQueue{{vertex: src, dist: 0}}
    heap.Init(pq)

    for pq.Len() > 0 {
        curr := heap.Pop(pq).(*pqItem)
        if curr.dist > dist[curr.vertex] {
            continue // stale entry
        }
        for _, e := range g.Adj[curr.vertex] {
            newDist := dist[curr.vertex] + e.Weight
            if newDist < dist[e.To] {
                dist[e.To] = newDist
                parent[e.To] = curr.vertex
                road[e.To] = e.Name
                heap.Push(pq, &pqItem{vertex: e.To, dist: newDist})
            }
        }
    }
    return &Result{Dist: dist, Parent: parent, Road: road}
}

// Path reconstructs the route from src to dst.
func (r *Result) Path(dst int) (vertices []int, roads []string) {
    if r.Dist[dst] == math.MaxInt64 {
        return nil, nil // no path
    }
    for v := dst; v != -1; v = r.Parent[v] {
        vertices = append(vertices, v)
        if r.Parent[v] != -1 {
            roads = append(roads, r.Road[v])
        }
    }
    // Reverse
    for i, j := 0, len(vertices)-1; i < j; i, j = i+1, j-1 {
        vertices[i], vertices[j] = vertices[j], vertices[i]
    }
    for i, j := 0, len(roads)-1; i < j; i, j = i+1, j-1 {
        roads[i], roads[j] = roads[j], roads[i]
    }
    return vertices, roads
}
```

---

## 4. Main — CLI

```go
// main.go
package main

import (
    "flag"
    "fmt"
    "math"
    "os"
    "strings"

    "routeplanner/algo"
    "routeplanner/graph"
)

func main() {
    mapFile := flag.String("map", "maps/cities.json", "path to the city map JSON file")
    from    := flag.String("from", "", "starting city")
    to      := flag.String("to", "", "destination city")
    all     := flag.Bool("all", false, "print all shortest paths from --from")
    flag.Parse()

    g, err := graph.LoadFromFile(*mapFile)
    if err != nil {
        fmt.Fprintln(os.Stderr, "load map:", err)
        os.Exit(1)
    }

    if *from == "" {
        fmt.Fprintln(os.Stderr, "usage: routeplanner --from <city> [--to <city>] [--all]")
        fmt.Fprintln(os.Stderr, "available cities:", strings.Join(g.Cities, ", "))
        os.Exit(1)
    }

    srcID, ok := g.Index[*from]
    if !ok {
        fmt.Fprintf(os.Stderr, "unknown city %q\n", *from)
        os.Exit(1)
    }

    result := algo.Dijkstra(g, srcID)

    if *all {
        fmt.Printf("Shortest distances from %s:\n", *from)
        for id, city := range g.Cities {
            if id == srcID { continue }
            if result.Dist[id] == math.MaxInt64 {
                fmt.Printf("  %-20s unreachable\n", city)
                continue
            }
            fmt.Printf("  %-20s %5d km\n", city, result.Dist[id])
        }
        return
    }

    if *to == "" {
        fmt.Fprintln(os.Stderr, "specify --to <city> or use --all")
        os.Exit(1)
    }

    dstID, ok := g.Index[*to]
    if !ok {
        fmt.Fprintf(os.Stderr, "unknown city %q\n", *to)
        os.Exit(1)
    }

    if result.Dist[dstID] == math.MaxInt64 {
        fmt.Printf("No route from %s to %s\n", *from, *to)
        os.Exit(1)
    }

    vertices, roads := result.Path(dstID)
    fmt.Printf("Route: %s → %s\n", *from, *to)
    fmt.Printf("Total: %d km\n\n", result.Dist[dstID])

    for i, road := range roads {
        from := g.Cities[vertices[i]]
        to   := g.Cities[vertices[i+1]]
        fmt.Printf("  %s → %s via %s\n", from, to, road)
    }
}
```

---

## 5. Running It

```bash
go mod init routeplanner
go run . --map maps/cities.json --from A --to E

# Route: A → E
# Total: 250 km
#
#   A → B via Highway 1
#   B → C via Main St
#   C → E via Coastal Hwy

go run . --map maps/cities.json --from A --all
# Shortest distances from A:
#   B                      100 km
#   C                      140 km
#   D                      200 km
#   E                      250 km
```

---

## 6. Adding More Cities to the Map

```go
// Procedurally generate a large map for benchmarking
// testdata/gen.go
package main

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "os"
)

func main() {
    n := 1000  // cities
    cities := make([]string, n)
    for i := range cities { cities[i] = fmt.Sprintf("City%04d", i) }
    
    type road struct {
        From, To, Name string
        KM             int `json:"km"`
    }
    var roads []road
    
    // Connect each city to 3-5 neighbors
    for i := 0; i < n; i++ {
        neighbors := 3 + rand.Intn(3)
        for k := 0; k < neighbors; k++ {
            j := rand.Intn(n)
            if j == i { continue }
            roads = append(roads, road{
                From: cities[i], To: cities[j],
                KM:   10 + rand.Intn(490),
                Name: fmt.Sprintf("Road_%d_%d", i, j),
            })
        }
    }
    
    data, _ := json.MarshalIndent(map[string]any{
        "cities": cities,
        "roads":  roads,
    }, "", "  ")
    os.WriteFile("maps/large.json", data, 0o644)
    fmt.Printf("Generated %d cities, %d roads\n", n, len(roads))
}
```

---

## 7. Key Concepts Applied

| Concept (Vol 3-4) | Where Used |
|-------------------|------------|
| Priority queue (min-heap) | Dijkstra's `container/heap` |
| Adjacency list | `graph.Adj [][]Edge` |
| Graph representation | JSON-loaded weighted graph |
| Shortest path | Dijkstra with early exit on stale entries |
| Path reconstruction | `parent[]` array, reverse traversal |
| Negative weight check | Return error on negative km in loader |

---

## Exercises

### Easy
1. Add a `--via <city>` flag: find the shortest route that must pass through an intermediate city. Hint: run Dijkstra twice (src → via, via → dst) and combine.
2. Add a road `weight` option: besides km, compute the route with minimum number of road segments (unweighted BFS) and compare it to the shortest km route.
3. Display the route as a simple ASCII map: print the city names in a column, connected by `→` arrows.

### Medium
4. Add **bidirectional Dijkstra**: run Dijkstra simultaneously from both src and dst. Terminate when the two searches meet. This roughly halves the search space for large maps.
5. Add a `--blocked <road>` flag to mark a road as closed. Remove that edge from the adjacency list before running Dijkstra. This simulates road closures.
6. Implement **A\* search**: add latitude/longitude to each city, compute the straight-line (haversine) distance as the heuristic. A\* should explore significantly fewer nodes than Dijkstra on a realistic map.

### Hard
7. Extend the planner to support **K-shortest paths** (Yen's algorithm): find the k-th best route that doesn't repeat vertices. Useful for showing route alternatives.
8. Build a **REST API** around the route planner: expose endpoints `GET /route?from=A&to=E` and `GET /distances?from=A`. Use chi router (Ch 60) and return JSON responses.
