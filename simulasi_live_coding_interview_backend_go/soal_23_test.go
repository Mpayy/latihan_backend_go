package main

// Soal 23 — Container With Most Water (Review Two-Pointer)

// Review konsep two-pointer (dari soal 5 & 15), tapi dengan cara pakai yang berbeda dari sebelumnya — ini penting supaya kamu paham two-pointer itu bukan cuma "satu trik", tapi bisa diterapkan dengan logika yang berbeda-beda.

// Diberikan height, sebuah slice angka yang mewakili tinggi dinding di posisi 0, 1, 2, ..., n-1. Pilih dua dinding yang, bersama sumbu-x di antaranya, membentuk sebuah wadah (container) yang bisa menampung air sebanyak-banyaknya. Return kapasitas air maksimum itu.

// go
// func MaxArea(height []int) int

// Contoh:

// Input:  height = [1,8,6,2,5,4,8,3,7]
// Output: 49
// Dinding di index 1 (tinggi 8) dan index 8 (tinggi 7)
// Lebar = 8-1 = 7, tinggi wadah = min(8,7) = 7
// Luas = 7 * 7 = 49

func MaxArea(height []int) int {
	left := 0
	right := len(height) - 1
	maxWater := 0

	for left < right {
		lebar := right - left
		tinggiWadah := min(height[left], height[right])

		currMax := lebar * tinggiWadah
		if currMax > maxWater {
			maxWater = currMax
		}

		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}

	return maxWater
}
