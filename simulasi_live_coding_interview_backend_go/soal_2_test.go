package main

import (
	"fmt"
	"slices"
	"testing"
)

// Soal 2
// "Group Anagrams"

// Diberikan sebuah slice of string, kelompokkan string-string yang merupakan anagram satu sama lain ke dalam grup yang sama.
// Dua string dianggap anagram kalau tersusun dari huruf yang sama, hanya urutannya berbeda (misalnya "eat" dan "tea" adalah anagram).

// Signature function:

// func GroupAnagrams(strs []string) [][]string

// Contoh:

// Input:  []string{"eat", "tea", "tan", "ate", "nat", "bat"}
// Output: [][]string{
//     {"eat", "tea", "ate"},
//     {"tan", "nat"},
//     {"bat"},
// }

// urutan grup dan urutan dalam grup tidak masalah

func GroupAnagrams(strs []string) [][]string {
	mapStr := make(map[string][]string, 0)

	for _, str := range strs {
		r := []rune(str)
		slices.Sort(r)
		sortStr := string(r)
		mapStr[sortStr] = append(mapStr[sortStr], str)
	}

	fmt.Println(mapStr)

	sliceStr := make([][]string, 0, len(mapStr))
	for _, vKey := range mapStr {
		sliceStr = append(sliceStr, vKey)
	}

	return sliceStr
}

// Helper function untuk mengurutkan hasil [][]string agar mudah dibandingkan
func normalizeResult(groups [][]string) {
	for _, group := range groups {
		slices.Sort(group) // Urutkan kata-kata di dalam setiap grup
	}
	// Urutkan grup berdasarkan kata pertamanya
	slices.SortFunc(groups, func(a, b []string) int {
		if len(a) == 0 || len(b) == 0 {
			return 0
		}
		if a[0] < b[0] {
			return -1
		} else if a[0] > b[0] {
			return 1
		}
		return 0
	})
}

func isSameGroups(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	normalizeResult(a)
	normalizeResult(b)

	for i := range a {
		if !slices.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestGroupAnagrams(t *testing.T) {
	strs := []string{"eat", "tea", "tan", "ate", "nat", "bat"}

	expected := [][]string{
		{"bat"},
		{"nat", "tan"},
		{"ate", "eat", "tea"},
	}

	res := GroupAnagrams(strs)

	if !isSameGroups(res, expected) {
		t.Errorf("Hasil salah!\nEkspektasi: %v\nDapatnya   : %v", expected, res)
	}
}

func TestGroupAnagrams_EmptyAndSingle(t *testing.T) {
	// Case 1: String kosong
	res1 := GroupAnagrams([]string{""})
	expected1 := [][]string{{""}}
	if !isSameGroups(res1, expected1) {
		t.Errorf("Case empty string salah! Ekspektasi %v, dapat %v", expected1, res1)
	}

	// Case 2: Satu karakter
	res2 := GroupAnagrams([]string{"a"})
	expected2 := [][]string{{"a"}}
	if !isSameGroups(res2, expected2) {
		t.Errorf("Case single char salah! Ekspektasi %v, dapat %v", expected2, res2)
	}
}
