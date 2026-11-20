package main

// Soal 42: Kampanye Tanpa Hari Berurutan

// Review DP dari soal 22 dengan kasus baru. Kamu punya perkiraan omzet harian dari sebuah kampanye promo. Kampanye tidak boleh dijalankan di dua hari yang berurutan (tim butuh jeda). Pilih hari-hari kampanye yang total omzetnya paling besar.

// go
// func MaxCampaignRevenue(revenue []int) int

// Contoh:

// [1,2,3,1]    → 4    // hari 1 dan 3 (1 + 3)
// [2,7,9,3,1]  → 12   // hari 1, 3, dan 5 (2 + 9 + 1)
// [5]          → 5
// []           → 0

// Requirement: O(N) time, dan O(1) space adalah nilai tambah. Format sama: satu pesan berisi clarifying question, pendekatan, kode, dan Big-O. Di pendekatanmu, jelaskan keputusan di tiap hari (ambil atau lewat) dan apa yang perlu diingat dari hari-hari sebelumnya.

func MaxCampaignRevenue(revenue []int) int {
	prev1 := 0
	prev2 := 0

	for _, rev := range revenue {
		curr := max(prev1, prev2+rev)
		prev2 = prev1
		prev1 = curr
	}
	return prev1
}
