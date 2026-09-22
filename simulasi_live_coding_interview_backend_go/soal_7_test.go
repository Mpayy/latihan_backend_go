package main

import (
	"cmp"
	"slices"
	"testing"
)

// Soal 7 — Merge Intervals
// Ini soal array/pattern klasik yang sering muncul juga di kasus nyata backend (misalnya: menggabungkan slot waktu meeting yang bertabrakan, menggabungkan range IP, dll).

// Diberikan sejumlah interval (rentang), masing-masing direpresentasikan sebagai [start, end]. Gabungkan semua interval yang saling tumpang tindih (overlapping) menjadi satu interval saja.

// func MergeIntervals(intervals [][]int) [][]int

// Contoh:
// Input:  [[1,3],[2,6],[8,10],[15,18]]
// Output: [[1,6],[8,10],[15,18]]
// [1,3] dan [2,6] overlap (karena 2 <= 3), digabung jadi [1,6]
// [8,10] dan [15,18] tidak overlap dengan yang lain, tetap terpisah

// Input:  [[1,4],[4,5]]
// Output: [[1,5]]
// [1,4] dan [4,5] dianggap overlap karena bersentuhan di angka 4

func MergeIntervals(intervals [][]int) [][]int {
	var result [][]int

	newIntervals := make([][]int, len(intervals))
	for i, v := range intervals {
		newIntervals[i] = slices.Clone(v)
	}

	slices.SortFunc(newIntervals, func(a, b []int) int {
		if n := cmp.Compare(a[0], b[0]); n != 0 {
			return n
		}

		return cmp.Compare(a[1], b[1])
	})

	for _, v := range newIntervals {
		// result kosong
		if len(result) == 0 {
			result = append(result, v)
			continue
		}

		// overlap
		lastIdx := len(result) - 1
		if v[0] <= result[lastIdx][1] {
			result[lastIdx][1] = max(result[lastIdx][1], v[1])
			continue
		}

		// not overlap
		result = append(result, v)
	}

	return result
}

// Helper function untuk membandingkan dua slice 2D [][]int
func equalIntervals(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !slices.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

func TestMergeIntervals(t *testing.T) {
	// Definisi kasus uji (Table-Driven Tests)
	tests := []struct {
		name      string
		intervals [][]int
		expected  [][]int
	}{
		{
			name:      "Kasus standar overlap bertahap",
			intervals: [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}},
			expected:  [][]int{{1, 6}, {8, 10}, {15, 18}},
		},
		{
			name:      "Kasus overlap sempurna/satu menutupi lainnya",
			intervals: [][]int{{1, 4}, {4, 5}},
			expected:  [][]int{{1, 5}},
		},
		{
			name:      "Kasus input tidak terurut",
			intervals: [][]int{{15, 18}, {1, 3}, {2, 6}, {8, 10}},
			expected:  [][]int{{1, 6}, {8, 10}, {15, 18}},
		},
		{
			name:      "Kasus interval saling mencakup (nested)",
			intervals: [][]int{{1, 10}, {2, 3}, {4, 8}},
			expected:  [][]int{{1, 10}},
		},
		{
			name:      "Kasus tidak ada overlap sama sekali",
			intervals: [][]int{{1, 2}, {3, 4}, {5, 6}},
			expected:  [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:      "Kasus satu interval",
			intervals: [][]int{{1, 4}},
			expected:  [][]int{{1, 4}},
		},
		{
			name:      "Kasus slice kosong",
			intervals: [][]int{},
			expected:  nil, // atau [][]int(nil)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := MergeIntervals(tt.intervals)

			if !equalIntervals(res, tt.expected) {
				t.Errorf("Hasil salah untuk %s!\nEkspektasi: %v\nDapatnya   : %v", tt.name, tt.expected, res)
			}
		})
	}
}
