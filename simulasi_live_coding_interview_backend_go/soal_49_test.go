package main

import (
	"fmt"
	"testing"
)

// Soal 49: Jumlah Kuadrat Minimum

// Sebuah sistem kuota menguraikan bilangan n menjadi jumlah kuadrat sempurna (1, 4, 9, 16, ...) dengan jumlah suku sesedikit mungkin. Suku yang sama boleh dipakai berkali-kali.

// go
// func MinSquares(n int) int

// Contoh:

// 12 → 3    // 4 + 4 + 4
// 13 → 2    // 4 + 9
// 1  → 1

// Hal yang tidak disebut (misalnya n nol atau negatif) sengaja tidak dijelaskan, itu bagian clarifying question-mu. Requirement: O(n × √n).

// Format satu pesan: clarifying question, pendekatan, kode, dan Big-O. Di pendekatanmu, jelaskan soal sebelumnya mana yang mirip dengan ini, dan apa yang berperan sebagai "daftar koin" di sini. Coba dulu tanpa membuka catatan soal 46 dan 48. Kalau buntu, tulis "hint 1".

func MinSquares(n int) int {
	switch {
	case n == 0:
		return 0
	case n < 0:
		return -1
	}

	dp := make([]int, n+1)

	for i := 1; i <= n; i++ {
		dp[i] = n + 1
		for j := 1; j*j <= i; j++ {
			dp[i] = min(dp[i], dp[i-j*j]+1)
		}
	}

	if dp[n] > n {
		return -1
	}

	return dp[n]
}

func TestMinSquares(t *testing.T) {
	res := MinSquares(12)
	fmt.Println(res)
}
