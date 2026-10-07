package main

import (
	"fmt"
	"testing"
)

// Soal 21 — Trie (Prefix Tree) untuk Autocomplete Sederhana

// Konsep yang perlu dipahami dulu kalau belum familiar: Trie (dibaca "try"), struktur data pohon yang dioptimalkan untuk pencarian berdasarkan prefix (awalan string) — sering dipakai untuk fitur autocomplete, spell-checker, atau validasi kata.

// Familiar dengan konsep Trie sebelumnya?

// Kalau belum, saya jelaskan dulu primernya sebelum soal. Kalau sudah pernah dengar/pakai, langsung saja ke soal:

// Soal: Implementasikan struktur data Trie dengan 2 operasi dasar:

// go
// type Trie struct {
//     // isi sendiri
// }

// func NewTrie() *Trie
// func (t *Trie) Insert(word string)
// func (t *Trie) Search(word string) bool  // true kalau `word` persis ada di trie (bukan cuma prefix)

// Contoh:

// go
// trie := NewTrie()
// trie.Insert("apple")
// trie.Search("apple")   // true
// trie.Search("app")     // false (cuma prefix, bukan kata yang di-insert)
// trie.Insert("app")
// trie.Search("app")     // true (sekarang sudah di-insert juga)

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}
}

func (t *Trie) Insert(word string) {
	curr := t.root

	for _, char := range word {
		node, exists := curr.children[char]
		if !exists {
			node = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
			curr.children[char] = node
		}

		curr = node
	}
	
	curr.isEnd = true
}

func (t *Trie) Search(word string) bool {
	curr := t.root

	for _, char := range word {
		node, exists := curr.children[char]
		if !exists {
			return false
		}

		curr = node
	}

	return curr.isEnd
}

func TestTrie(t *testing.T) {
	trie := NewTrie()
	trie.Insert("apple")
	fmt.Println(trie.Search("apple"))
	fmt.Println(trie.Search("app"))
	trie.Insert("app")
	fmt.Println(trie.Search("app"))
}
