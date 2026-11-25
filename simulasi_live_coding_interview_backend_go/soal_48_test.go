package main

import (
	"fmt"
	"testing"
)

// Soal 48: Kardus Minimum

// Ini soal 46 dengan studi kasus berbeda, supaya polanya makin melekat. Kalau perlu, buka catatan tabel dp dari penjelasan soal 46, itu boleh.

// Sebuah gudang mengirim barang dalam kardus berukuran tetap. Tersedia beberapa ukuran kardus (jumlah barang yang muat), masing-masing dalam jumlah tak terbatas. Untuk pesanan items barang, pakai kardus sesedikit mungkin, dan setiap kardus harus terisi penuh.

// go
// func MinBoxes(sizes []int, items int) int

// Return jumlah kardus minimum, atau -1 kalau pesanan tidak bisa dibentuk.

// Contoh:

// sizes=[5,3,4], items=7  → 2    // 3 + 4
// sizes=[6,4,3], items=12 → 2    // 6 + 6
// sizes=[5],     items=7  → -1

// Hal yang tidak disebut (misalnya items nol, daftar ukuran kosong, atau ukuran nol atau negatif) sengaja tidak dijelaskan, itu bagian clarifying question-mu. Requirement: O(items × jumlah ukuran).

// Format satu pesan: clarifying question, pendekatan, kode, dan Big-O. Di pendekatanmu, jelaskan:

// Kenapa memakai kardus terbesar dulu tidak cukup, dengan trace contoh pertama.
// Arti dp[i], nilai dasarnya, dan apa yang kamu isikan untuk nominal yang tidak bisa dibentuk.

// Kalau buntu, tulis "hint 1".

func MinBoxes(sizes []int, items int) int {
	switch {
	case items == 0:
		return 0
	case items < 0:
		return -1
	}

	dp := make([]int, items+1)

	for i := 1; i <= items; i++ {
		dp[i] = items + 1
	}

	for i := 1; i <= items; i++ {
		for _, s := range sizes {
			if s > 0 && s <= i && dp[i-s] <= items {
				dp[i] = min(dp[i], dp[i-s]+1)
			}
		}
	}

	if dp[items] > items {
		return -1
	}

	return dp[items]
}

func TestMinBoxes(t *testing.T) {
	sizes := []int{5, 3, 4}
	items := 7

	res := MinBoxes(sizes, items)
	fmt.Println(res)
}
