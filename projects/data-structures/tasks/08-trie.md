---
title: Trie
step: 8
difficulty: medium
estimated: 40 min
---

## What You Are Building

A trie (rhymes with "try") is a tree for storing strings character by character. Every path from root to a marked node spells out a word. Tries power autocomplete, spell checkers, IP routing tables, and word games.

```
insert("go", "got", "get", "dog")

        root
       /    \
      g      d
     / \      \
    o   e      o
   /|    \      \
  (*)t   t      g
      \    \     \
      (*)  (*)   (*)

(*) = isEnd=true
```

## Key Concepts

**TrieNode structure** — Each node stores:
- `children map[rune]*TrieNode`: which characters can follow this one
- `isEnd bool`: does a valid word end at this node

Using `map[rune]` instead of a fixed `[26]byte` array handles Unicode and makes the code clearer.

**Insert** — Walk character by character. If a child for the current rune doesn't exist, create it. After processing all characters, mark the final node as `isEnd = true`.

**Search** — Walk character by character. If any character is missing from the children, the word doesn't exist. After the last character, return `node.isEnd` — this distinguishes "go" (complete word) from "g" (prefix only).

**StartsWith** — Same as Search but return `true` if you reach the end of the prefix without any missing character, regardless of `isEnd`.

**WordsWithPrefix** — Start from the prefix's ending node, then DFS to collect all paths that end at `isEnd` nodes.

## Struct Signatures

```go
type TrieNode struct {
    children map[rune]*TrieNode
    isEnd    bool
}

type Trie struct {
    root *TrieNode
}

func NewTrie() *Trie {
    return &Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}}
}
```

## Methods to Implement

| Method | Description |
|--------|-------------|
| `Insert(word string)` | Add word to the trie |
| `Search(word string) bool` | True if the exact word was inserted |
| `StartsWith(prefix string) bool` | True if any inserted word starts with prefix |
| `WordsWithPrefix(prefix string) []string` | All inserted words starting with prefix |

## Edge Cases to Handle

- `Search("")`: return `true` if empty string was inserted, `false` otherwise
- `StartsWith("")`: return `true` (every word starts with the empty prefix)
- `WordsWithPrefix` for a prefix that matches no words: return `[]string{}`
- Inserting the same word twice: no duplicate in results

## Example

```go
t := NewTrie()
for _, w := range []string{"apple", "app", "application", "apply", "apt"} {
    t.Insert(w)
}

fmt.Println(t.Search("app"))         // true
fmt.Println(t.Search("ap"))          // false (not inserted as a whole word)
fmt.Println(t.StartsWith("ap"))      // true
fmt.Println(t.StartsWith("xyz"))     // false

words := t.WordsWithPrefix("app")
sort.Strings(words)
fmt.Println(words) // [app apple application apply]
```

## Hints

- Write a private helper `func (t *Trie) nodeAt(s string) *TrieNode` that walks to the node at the end of `s`, returning `nil` if any character is missing. Both `Search` and `StartsWith` can use it.
- `WordsWithPrefix`: first find the node at the prefix endpoint, then run a DFS from that node accumulating the prefix + characters explored.
- For DFS in `WordsWithPrefix`, write a recursive helper: `func collect(node *TrieNode, current string, results *[]string)`.
