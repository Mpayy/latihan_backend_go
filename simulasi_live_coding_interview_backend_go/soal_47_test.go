package main

import (
	"fmt"
	"testing"
)

// Soal 47: Tangga dengan Biaya

// Setiap anak tangga punya biaya yang dibayar saat kamu menginjaknya. Dari anak tangga i, kamu bisa naik 1 atau 2 anak tangga. Kamu boleh mulai dari index 0 atau 1. Puncak adalah posisi setelah anak tangga terakhir. Cari total biaya minimum untuk mencapai puncak.

// go
// func MinCostClimb(cost []int) int

// Contoh:

// [10,15,20]                    → 15   // mulai di index 1, bayar 15, lompat 2 langkah ke puncak
// [1,100,1,1,1,100,1,1,100,1]   → 6

// Hal yang tidak disebut (misalnya panjang minimum atau biaya nol atau negatif) sengaja tidak dijelaskan, itu bagian clarifying question-mu. Format satu pesan: clarifying question, pendekatan, kode, dan Big-O. Di pendekatanmu, jelaskan arti dp[i], dari mana nilainya dihitung, dan apakah bisa O(1) space. Kalau buntu, tulis "hint 1".

func MinCostClimb(cost []int) int {
	if len(cost) < 2 {
		return 0
	}

	prev1 := 0
	prev2 := 0

	for i := 2; i <= len(cost); i++ {
		curr := min(prev1+cost[i-1], prev2+cost[i-2])
		prev2 = prev1
		prev1 = curr
	}

	return prev1
}

// func MinCostClimb(cost []int) int {
// 	if len(cost) < 2 {
// 		return 0
// 	}

// 	dp := make([]int, len(cost)+1)

// 	for i := 2; i <= len(cost); i++ {
// 		dp[i] = min(dp[i-1]+cost[i-1], dp[i-2]+cost[i-2])
// 	}

// 	return dp[len(cost)]
// }

func TestMinCostClimb(t *testing.T) {
	cost := []int{1, 100, 1, 1, 1, 100, 1, 1, 100, 1}
	res := MinCostClimb(cost)
	fmt.Println(res)
}
