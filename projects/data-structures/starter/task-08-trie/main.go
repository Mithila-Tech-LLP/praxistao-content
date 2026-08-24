package main

import "sort"

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

func (t *Trie) Insert(word string) {
	// TODO: implement
}

func (t *Trie) Search(word string) bool {
	// TODO: implement
	return false
}

func (t *Trie) StartsWith(prefix string) bool {
	// TODO: implement
	return false
}

func (t *Trie) WordsWithPrefix(prefix string) []string {
	// TODO: implement — return all words that start with prefix (sorted)
	return []string{}
}

var _ = sort.Strings // hint
