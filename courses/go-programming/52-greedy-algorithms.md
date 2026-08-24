# Chapter 37: Greedy Algorithms

A greedy algorithm makes the locally optimal choice at each step, hoping to find a globally optimal solution. Unlike DP, it never reconsiders past choices. This makes greedy algorithms fast and simple — when they work. Knowing *when* to trust greedy intuition is the skill.

## Table of Contents

1. [The Greedy Approach](#1-the-greedy-approach)
2. [Activity Selection and Interval Problems](#2-activity-selection-and-interval-problems)
3. [Scheduling and Ordering](#3-scheduling-and-ordering)
4. [Graph Greedy — MST and Dijkstra](#4-graph-greedy--mst-and-dijkstra)
5. [Array and String Greedy](#5-array-and-string-greedy)
6. [Proving Greedy Correctness](#6-proving-greedy-correctness)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. The Greedy Approach

**When greedy works:**
- Greedy choice property: a globally optimal solution can be built by always choosing the locally optimal (greedy) option
- Optimal substructure: an optimal solution contains optimal solutions to subproblems

**Greedy vs DP:**
- Coin change with denominations {1, 5, 10, 25}: greedy (always pick largest ≤ remaining) works
- Coin change with denominations {1, 3, 4}: greedy fails for target 6 (gives 4+1+1=3 coins; optimal is 3+3=2 coins) → need DP

**The proof technique**: exchange argument — assume the optimal solution differs from the greedy solution; show that swapping the optimal choice for the greedy choice doesn't make it worse.

---

## 2. Activity Selection and Interval Problems

The canonical greedy problem: select maximum non-overlapping activities.

```go
type Interval struct{ Start, End int }

// ActivitySelection: select max number of non-overlapping activities.
// Greedy: sort by end time, always pick the activity that finishes earliest.
func ActivitySelection(activities []Interval) []Interval {
    if len(activities) == 0 { return nil }
    sort.Slice(activities, func(i, j int) bool {
        return activities[i].End < activities[j].End
    })

    selected := []Interval{activities[0]}
    lastEnd := activities[0].End

    for _, act := range activities[1:] {
        if act.Start >= lastEnd {  // No overlap
            selected = append(selected, act)
            lastEnd = act.End
        }
    }
    return selected
}
// Why sort by end? Finishing earliest leaves maximum room for future activities.
```

### Merge Overlapping Intervals
```go
// MergeIntervals: merge all overlapping intervals.
func MergeIntervals(intervals []Interval) []Interval {
    if len(intervals) == 0 { return nil }
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i].Start < intervals[j].Start
    })

    merged := []Interval{intervals[0]}
    for _, curr := range intervals[1:] {
        last := &merged[len(merged)-1]
        if curr.Start <= last.End {  // Overlap
            if curr.End > last.End { last.End = curr.End }
        } else {
            merged = append(merged, curr)
        }
    }
    return merged
}
```

### Minimum Rooms for Meeting Schedules
```go
// MinMeetingRooms: minimum rooms needed to schedule all meetings.
// Greedy: track when meetings start and end using separate sorted arrays.
func MinMeetingRooms(intervals []Interval) int {
    starts := make([]int, len(intervals))
    ends := make([]int, len(intervals))
    for i, iv := range intervals {
        starts[i] = iv.Start
        ends[i] = iv.End
    }
    sort.Ints(starts)
    sort.Ints(ends)

    rooms, endPtr := 0, 0
    for _, s := range starts {
        if s < ends[endPtr] {
            rooms++  // New room needed
        } else {
            endPtr++  // Reuse a room that just freed
        }
    }
    return rooms
}
```

### Jump Game
```go
// CanJump: given jump lengths, can you reach the last index?
// Greedy: track the farthest index reachable so far.
func CanJump(nums []int) bool {
    maxReach := 0
    for i, jump := range nums {
        if i > maxReach { return false }  // Can't reach index i
        if reach := i + jump; reach > maxReach {
            maxReach = reach
        }
    }
    return true
}

// JumpGame2: minimum jumps to reach end.
func MinJumps(nums []int) int {
    jumps, currentEnd, farthest := 0, 0, 0
    for i := 0; i < len(nums)-1; i++ {
        if reach := i + nums[i]; reach > farthest {
            farthest = reach
        }
        if i == currentEnd {  // Exhausted current jump range
            jumps++
            currentEnd = farthest
        }
    }
    return jumps
}
```

---

## 3. Scheduling and Ordering

### Fractional Knapsack (vs 0/1 Knapsack)
```go
type Item struct{ Weight, Value float64 }

// FractionalKnapsack: can take fractions of items. Greedy by value/weight ratio.
func FractionalKnapsack(items []Item, capacity float64) float64 {
    sort.Slice(items, func(i, j int) bool {
        ri := items[i].Value / items[i].Weight
        rj := items[j].Value / items[j].Weight
        return ri > rj  // Highest ratio first
    })

    totalValue := 0.0
    for _, item := range items {
        if capacity <= 0 { break }
        take := min64(item.Weight, capacity)
        totalValue += take * (item.Value / item.Weight)
        capacity -= take
    }
    return totalValue
}

func min64(a, b float64) float64 {
    if a < b { return a }
    return b
}
// Fractional knapsack: greedy works. 0/1 knapsack: greedy fails → need DP.
```

### Task Scheduler
```go
// TaskScheduler: given tasks and cooldown n, minimum time to complete all.
// Greedy: always schedule the most frequent remaining task.
func LeastInterval(tasks []byte, n int) int {
    freq := [26]int{}
    for _, t := range tasks { freq[t-'A']++ }
    sort.Ints(freq[:])

    maxFreq := freq[25]
    // Slots: (maxFreq-1) chunks of (n+1), plus the last chunk of tasks with max freq
    idleSlots := (maxFreq - 1) * (n + 1)

    // Count how many tasks have max frequency (they go in the last chunk):
    lastChunk := 0
    for _, f := range freq {
        if f == maxFreq { lastChunk++ }
    }

    total := idleSlots + lastChunk
    if total < len(tasks) { return len(tasks) }
    return total
}
```

### Gas Station
```go
// CanCompleteCircuit: find starting gas station where you can complete the circuit.
// Greedy: if total gas >= total cost, a solution exists. Start from station after deficit.
func CanCompleteCircuit(gas, cost []int) int {
    totalSurplus, currentSurplus, start := 0, 0, 0

    for i := range gas {
        diff := gas[i] - cost[i]
        totalSurplus += diff
        currentSurplus += diff
        if currentSurplus < 0 {
            start = i + 1  // Can't start from anywhere up to i
            currentSurplus = 0
        }
    }
    if totalSurplus < 0 { return -1 }
    return start
}
```

---

## 4. Graph Greedy — MST and Dijkstra

### Prim's MST
```go
// Prim's: build MST by always adding the cheapest edge connecting the MST to a new vertex.
func PrimMST(n int, edges [][]int) int {  // edges: [u, v, weight]
    adj := make([][][2]int, n)  // adj[u] = [{v, weight}]
    for _, e := range edges {
        adj[e[0]] = append(adj[e[0]], [2]int{e[1], e[2]})
        adj[e[1]] = append(adj[e[1]], [2]int{e[0], e[2]})
    }

    inMST := make([]bool, n)
    key := make([]int, n)
    for i := range key { key[i] = 1<<62 }
    key[0] = 0

    pq := &MinPQ{}  // Min priority queue on (key, vertex)
    heap.Init(pq)
    heap.Push(pq, [2]int{0, 0})

    totalCost := 0
    for pq.Len() > 0 {
        item := heap.Pop(pq).([2]int)
        cost, u := item[0], item[1]
        if inMST[u] { continue }
        inMST[u] = true
        totalCost += cost

        for _, neighbor := range adj[u] {
            v, w := neighbor[0], neighbor[1]
            if !inMST[v] && w < key[v] {
                key[v] = w
                heap.Push(pq, [2]int{w, v})
            }
        }
    }
    return totalCost
}
```

### Kruskal's MST (Union-Find)
```go
// Kruskal's: sort all edges by weight; add edge if it doesn't create a cycle.
type Edge struct{ U, V, W int }

func KruskalMST(n int, edges []Edge) int {
    sort.Slice(edges, func(i, j int) bool { return edges[i].W < edges[j].W })

    parent := make([]int, n)
    rank := make([]int, n)
    for i := range parent { parent[i] = i }

    var find func(int) int
    find = func(x int) int {
        if parent[x] != x { parent[x] = find(parent[x]) }
        return parent[x]
    }

    union := func(x, y int) bool {
        px, py := find(x), find(y)
        if px == py { return false }  // Already connected — adding would create cycle
        if rank[px] < rank[py] { px, py = py, px }
        parent[py] = px
        if rank[px] == rank[py] { rank[px]++ }
        return true
    }

    totalCost, edgesUsed := 0, 0
    for _, e := range edges {
        if union(e.U, e.V) {
            totalCost += e.W
            edgesUsed++
            if edgesUsed == n-1 { break }
        }
    }
    return totalCost
}
```

---

## 5. Array and String Greedy

### Maximum Subarray (Kadane's)
```go
// MaxSubarray: maximum sum contiguous subarray.
// Greedy: extend current subarray if it adds value; restart otherwise.
func MaxSubarray(nums []int) int {
    maxSum := nums[0]
    curSum := nums[0]

    for _, n := range nums[1:] {
        if curSum < 0 { curSum = 0 }  // Starting fresh is better
        curSum += n
        if curSum > maxSum { maxSum = curSum }
    }
    return maxSum
}
```

### Non-overlapping Intervals (Minimum Removals)
```go
// EraseOverlapIntervals: minimum number of intervals to remove so none overlap.
// Equivalent to: (total intervals) - (max non-overlapping intervals)
// Greedy: sort by end, keep the most (same as activity selection).
func EraseOverlapIntervals(intervals []Interval) int {
    if len(intervals) == 0 { return 0 }
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i].End < intervals[j].End
    })

    keep, lastEnd := 1, intervals[0].End
    for _, iv := range intervals[1:] {
        if iv.Start >= lastEnd {
            keep++
            lastEnd = iv.End
        }
    }
    return len(intervals) - keep
}
```

### Assign Cookies
```go
// AssignCookies: maximize content children. Cookie s satisfies child g if s >= g.
// Greedy: sort both, assign smallest sufficient cookie to each child.
func FindContentChildren(g, s []int) int {
    sort.Ints(g)
    sort.Ints(s)

    child, cookie := 0, 0
    for child < len(g) && cookie < len(s) {
        if s[cookie] >= g[child] { child++ }  // Satisfied — move to next child
        cookie++  // Always move to next cookie
    }
    return child
}
```

### Partition Labels
```go
// PartitionLabels: partition string so each letter appears in at most one part.
// Greedy: extend current partition to include last occurrence of every seen letter.
func PartitionLabels(s string) []int {
    last := [26]int{}
    for i, c := range s { last[c-'a'] = i }

    sizes := []int{}
    start, end := 0, 0
    for i, c := range s {
        if last[c-'a'] > end { end = last[c-'a'] }
        if i == end {
            sizes = append(sizes, end-start+1)
            start = end + 1
        }
    }
    return sizes
}
```

---

## 6. Proving Greedy Correctness

### Exchange Argument Template
To prove that greedy solution G is optimal:
1. Assume the optimal solution O differs from G at some step
2. Show you can "exchange" O's choice for G's choice without decreasing the objective
3. Repeat until O = G — contradiction (O was supposed to be strictly better)

**Example: Activity Selection proof sketch**
- G sorts by earliest end time; O might pick activity `a` instead of `b` where `b.End < a.End`
- Swap `a` for `b` in O: `b.End ≤ a.End` means everything after still fits
- O is still valid and has the same number of activities
- So earliest-end is always at least as good → greedy is optimal

### Common Greedy Failure Modes
```
Problem                          | Why greedy fails
---------------------------------|------------------------------------------
0/1 Knapsack                    | Taking high ratio item blocks better combinations
Shortest Path (neg weights)     | Need Bellman-Ford; greedy freezes committed edges
Coin change (arbitrary denoms)  | Larger coin can prevent finding exact change
Longest Path in general graph   | Greedy can follow dead ends
TSP                             | Local choice doesn't account for full circuit
```

---

## Summary

- Greedy works when **greedy choice property** holds: local optima lead to global optimum
- **Activity selection** → sort by end time; **Interval scheduling** uses the same insight
- **Fractional knapsack** → sort by value/weight ratio (greedy works); **0/1 knapsack** → DP
- **Kadane's algorithm** is greedy DP: reset when sum is negative
- **Kruskal / Prim** are greedy MST algorithms; **Dijkstra** is greedy shortest path (non-negative weights only)
- **Proving greedy**: use the exchange argument — show swapping to greedy doesn't worsen the solution
- **Spotting greedy problems**: "maximize selections", "minimize cost/time", problems with natural ordering/sorting

---

## Exercises

### Easy
1. Implement `LemonadeChange(bills []int) bool`. Customers pay with 5, 10, or 20 dollar bills; lemonade costs $5. Make correct change for each customer or return false. Greedy: always break 20s with one 10 + one 5 before using three 5s (preserve 5s for 10s).
2. Implement `MonotonicArray(arr []int) bool` using a single pass. Then implement `MakeMonotonicMinRemovals(arr []int) int` — minimum removals to make the array non-decreasing. This is `len(arr) - LIS(arr)`.
3. Sort `[["Alice",40],["Bob",50],["Carol",30]]` by score descending using `sort.Slice`. Then implement a stable version where ties preserve alphabetical order using `sort.SliceStable`.

### Medium
4. **Two City Scheduling**: `2n` people need to fly to city A or city B. `costs[i] = [costA, costB]`. Exactly `n` fly to each city. Minimize total cost. Greedy: sort by `costA - costB` (biggest "prefer A" first → send first n to A, rest to B). Verify: `[[10,20],[30,200],[400,50],[30,20]]` → 110.
5. **Minimum Number of Arrows to Burst Balloons**: Balloons on x-axis as `[start, end]`. One arrow shot vertically at x pops all balloons where `start ≤ x ≤ end`. Minimize arrows needed. Greedy: sort by end; one arrow per non-overlapping group. Verify: `[[10,16],[2,8],[1,6],[7,12]]` → 2.
6. **Reorganize String**: rearrange characters so no two adjacent are the same. Greedy: always place the most frequent remaining character (that isn't the same as the last placed). Return `""` if impossible. Use a max-heap. Verify: `"aab"` → `"aba"`.

### Hard
7. **Candy**: n children in a row with ratings. Each child gets at least 1 candy. Higher-rated child must get more candy than their neighbor. Minimize total candies. Two-pass greedy: (1) left-to-right: if rating[i] > rating[i-1], candy[i] = candy[i-1]+1; (2) right-to-left: if rating[i] > rating[i+1], candy[i] = max(candy[i], candy[i+1]+1). Prove both passes are necessary. Verify: `[1,0,2]` → 5, `[1,2,2]` → 4.
8. **Largest Number**: Given integers, arrange them to form the largest possible number. Greedy: custom comparator where `a < b` if `str(a)+str(b) > str(b)+str(a)`. This is a non-obvious ordering — prove it forms a valid total order (transitivity). Verify: `[3,30,34,5,9]` → `"9534330"`.
