package main

import (
	"fmt"
	"testing"
)

// Soal 22 — Climbing Stairs (Intro ke Dynamic Programming)

// Konsep baru: Dynamic Programming (DP) — teknik menyelesaikan masalah dengan memecahnya jadi sub-masalah lebih kecil, dan menyimpan hasil sub-masalah itu supaya tidak dihitung ulang berkali-kali (teknik ini disebut memoization).

// Soal klasik paling umum untuk pengenalan DP:

// Kamu sedang menaiki tangga dengan n anak tangga. Setiap langkah, kamu bisa naik 1 atau 2 anak tangga. Ada berapa cara berbeda untuk mencapai anak tangga paling atas?

// go
// func ClimbStairs(n int) int

// Contoh:

// Input: n = 2
// Output: 2
// // Cara: (1 langkah + 1 langkah), (2 langkah)

// Input: n = 3
// Output: 3
// // Cara: (1+1+1), (1+2), (2+1)

// func ClimbStairs(n int) int {
// 	if n <= 1 {
// 		return 1
// 	}

// 	dp := make([]int, n+1)

// 	dp[1] = 1
// 	dp[2] = 2

// 	for i := 3; i <= n; i++ {
// 		dp[i] = dp[i-1] + dp[i-2]
// 	}

// 	return dp[n]
// }

func ClimbStairs(n int) int {
	prev1 := 1
	prev2 := 1

	for i := 2; i <= n; i++ {
		curr := prev1 + prev2
		prev2 = prev1
		prev1 = curr
	}

	return prev1
}

func TestClimbStairs(t *testing.T) {
	n := 5
	result := ClimbStairs(n)
	fmt.Println(result)
}
