package main

import (
	"fmt"
	"testing"
)

// Soal 46: Uang Kembalian Minimum

// Review DP (soal 22 dan 42). Sebuah kasir otomatis punya beberapa pecahan koin dalam jumlah tak terbatas. Cari jumlah koin paling sedikit untuk membentuk nominal tertentu.

// go
// func MinCoins(coins []int, amount int) int

// Return jumlah koin minimum, atau -1 kalau nominalnya tidak bisa dibentuk.

// Contoh:

// coins=[1,3,4],  amount=6  → 2    // 3 + 3
// coins=[1,2,5],  amount=11 → 3    // 5 + 5 + 1
// coins=[2],      amount=3  → -1

// Requirement: waktu O(amount × jumlah koin). Hal yang tidak disebut (misalnya amount nol, daftar koin kosong, atau nilai koin nol atau negatif) sengaja tidak dijelaskan, itu bagian clarifying question-mu.

// Format tetap satu pesan: clarifying question, pendekatan, kode, dan Big-O. Di pendekatanmu, jelaskan:

// Kenapa mengambil koin terbesar dulu tidak cukup, dengan trace contoh pertama.
// Apa arti dp[i], dan dari mana nilainya dihitung.

// Kalau terasa buntu, tulis "hint 1".

func MinCoins(coins []int, amount int) int {
	if amount < 0 {
		return -1
	}
	if amount == 0 {
		return 0
	}

	inf := amount + 1
	dp := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		dp[i] = inf
	}

	for i := 1; i <= amount; i++ {
		for _, c := range coins {
			if c > 0 && c <= i && dp[i-c] <= amount {
				dp[i] = min(dp[i], dp[i-c]+1)
			}
		}
	}

	if dp[amount] > amount {
		return -1
	}
	fmt.Println(dp)
	return dp[amount]
}

func TestMinCoins(t *testing.T) {
	coins := []int{1, 3, 4}
	amount := 6

	res := MinCoins(coins, amount)
	fmt.Println(res)
}
