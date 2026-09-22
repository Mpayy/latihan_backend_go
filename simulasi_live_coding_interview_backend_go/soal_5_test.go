package main

import "testing"

// Soal 5
// "Longest Substring Without Repeating Characters"

// Diberikan sebuah string s, cari panjang substring terpanjang yang tidak mengandung karakter berulang.

// func LengthOfLongestSubstring(s string) int

// Contoh:
// Input:  "abcabcbb"
// Output: 3   // substring-nya "abc"

// Input:  "bbbbb"
// Output: 1   // substring-nya "b"

// Input:  "pwwkew"
// Output: 3   // substring-nya "wke"

func LengthOfLongestSubstring(s string) int {
	lastSeen := make(map[rune]int)

	maxLength := 0
	left := 0

	for right, char := range s {
		if lastPos, exists := lastSeen[char]; exists && lastPos >= left {
			left = lastPos + 1
		}

		lastSeen[char] = right

		maxLength = max(maxLength, right-left+1)
	}
	return maxLength
}

func TestLengthOfLongestSubstring(t *testing.T) {
	// Definisi kasus uji (Table-Driven Tests)
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Kasus standar LeetCode 1",
			input:    "abcabcbb",
			expected: 3, // Substring: "abc"
		},
		{
			name:     "Kasus karakter sama semua",
			input:    "bbbbb",
			expected: 1, // Substring: "b"
		},
		{
			name:     "Kasus karakter unik berurutan",
			input:    "pwwkew",
			expected: 3, // Substring: "wke"
		},
		{
			name:     "Kasus string kosong",
			input:    "",
			expected: 0,
		},
		{
			name:     "Kasus satu karakter",
			input:    "a",
			expected: 1,
		},
		{
			name:     "Kasus karakter berulang di ujung",
			input:    "au",
			expected: 2,
		},
		{
			name:     "Kasus karakter berulang berdekatan",
			input:    "abba",
			expected: 2, // Substring: "ab" atau "ba"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := LengthOfLongestSubstring(tt.input)
			if res != tt.expected {
				t.Errorf("Hasil salah untuk input %q! Ekspektasi %d, tapi dapatnya %d", tt.input, tt.expected, res)
			}
		})
	}
}
