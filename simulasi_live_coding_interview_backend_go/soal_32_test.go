package main

import (
	"fmt"
	"testing"
)

// Soal 32 (Pemanasan): Duplikat yang Berdekatan

// Diberikan slice angka nums dan sebuah bilangan k. Return true kalau ada dua index berbeda i dan j sehingga nums[i] == nums[j] dan selisih index-nya paling banyak k (|i - j| <= k). Kalau tidak ada, return false.

// func ContainsNearbyDuplicate(nums []int, k int) bool

// Contoh:

// Input:  nums = [1,2,3,1], k = 3
// Output: true    // angka 1 ada di index 0 dan 3, selisihnya 3

// Input:  nums = [1,0,1,1], k = 1
// Output: true    // angka 1 di index 2 dan 3, selisihnya 1

// Input:  nums = [1,2,3,1,2,3], k = 2
// Output: false   // angka yang sama selalu berjarak 3

func ContainsNearbyDuplicate(nums []int, k int) bool {
	angkaIndex := make(map[int]int, 0)

	for i, num := range nums {
		index, exists := angkaIndex[num]
		if exists {
			if i-index <= k {
				return true
			}
		}

		angkaIndex[num] = i
	}

	return false
}

func TestContainsNearbyDuplicate(t *testing.T) {
	nums := []int{1,0,1,1}
	k := 1

	result := ContainsNearbyDuplicate(nums, k)
	fmt.Println(result)
}

