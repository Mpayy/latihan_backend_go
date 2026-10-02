package main

import (
	"fmt"
	"testing"
)

// Soal 13 — Binary Search Dasar (Review Soal 11)
// Studi kasus baru, tanpa rotasi dulu — murni untuk memastikan fondasi binary search-mu solid sebagai dasar sebelum nanti ketemu variasi yang lebih rumit lagi di sesi lain.

// Diberikan slice angka yang sudah terurut menaik (tanpa rotasi), dan sebuah target. Kalau target ada di slice, return index-nya. Kalau tidak ada, return index di mana target seharusnya disisipkan supaya slice tetap terurut.

// go
// func SearchInsert(nums []int, target int) int

// Contoh:

// Input:  nums = [1,3,5,6], target = 5
// Output: 2   // 5 ada di index 2

// Input:  nums = [1,3,5,6], target = 2
// Output: 1   // 2 tidak ada, tapi seharusnya disisipkan di index 1 (antara 1 dan 3)

// Input:  nums = [1,3,5,6], target = 7
// Output: 4   // 7 lebih besar dari semua elemen, disisipkan di akhir

func SearchInsert(nums []int, target int) int {
	left := 0
	right := len(nums) - 1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == target {
			return mid
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return left
}

func TestSearchInsert(t *testing.T) {
	nums := []int{1, 3, 5, 6}
	target := 7

	result := SearchInsert(nums, target)
	fmt.Println("Hasil return:", result)
}
