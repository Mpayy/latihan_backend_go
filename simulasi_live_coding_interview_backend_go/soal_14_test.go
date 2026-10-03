package main

import (
	"fmt"
	"testing"
)

// Soal 14 — First Unique Character
// Diberikan sebuah string s, cari index dari karakter pertama yang tidak berulang (muncul cuma sekali) di string itu. Kalau tidak ada karakter yang unik, return -1.

// go
// func FirstUniqChar(s string) int

// Contoh:
// Input:  "leetcode"
// Output: 0   // 'l' muncul cuma sekali, dan dia yang pertama

// Input:  "loveleetcode"
// Output: 2   // 'l' dan 'o' masing-masing muncul 2x, 'v' muncul sekali dan dia yang pertama kali muncul sebagai karakter unik

func FirstUniqChar(s string) int {
	mapStr := make(map[rune]int, 0)

	for _, char := range s {
		mapStr[char]++
	}

	for i, char := range s {
		if mapStr[char] == 1 {
			return i
		}
	}

	return -1
}

func TestFirstUniqChar(t *testing.T) {
	input := "loveleetcodev"
	result := FirstUniqChar(input)
	fmt.Println(result)
}
