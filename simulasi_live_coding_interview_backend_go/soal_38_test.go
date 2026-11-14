package main

// Soal 38: Hitung Rentang dengan Total Tertentu

// Ini variasi dari soal 25, tapi angkanya sekarang boleh negatif. Bayangkan mutasi stok harian yang bisa naik dan turun, dan kamu ingin tahu berapa rentang hari berurutan yang totalnya persis k.

// go
// func SubarraySum(nums []int, k int) int

// Hitung berapa banyak subarray berurutan yang jumlahnya persis k.

// Contoh:

// [1,1,1], k=2    → 2   // [1,1] di dua posisi
// [1,2,3], k=3    → 2   // [1,2] dan [3]
// [1,-1,0], k=0   → 3   // [1,-1], [1,-1,0], dan [0]

// Requirement: O(N). Format sama seperti tadi: kirim satu pesan berisi clarifying question, pendekatan, kode, dan Big-O. Di pendekatanmu, jelaskan kenapa sliding window dari soal 25 tidak cukup di sini, dengan trace singkat pakai contoh ketiga.

func SubarraySum(nums []int, k int) int {
	freq := make(map[int]int)
	freq[0] = 1
	currentSum := 0
	counter := 0

	for _, num := range nums {
		currentSum += num
		if count, exist := freq[currentSum-k]; exist {
			counter += count
		}
		freq[currentSum]++
	}

	return counter
}
