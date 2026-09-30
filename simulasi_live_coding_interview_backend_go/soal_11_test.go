package main

import (
	"fmt"
	"testing"
)

// Soal 11 — Binary Search on Rotated Sorted Array
// Ini pattern array yang sedikit lebih menantang dari soal-soal sebelumnya, tapi masih dalam kategori yang umum ditanyakan untuk posisi backend (variasi dari binary search, yang sering muncul di konteks nyata seperti pencarian data terurut).

// Diberikan sebuah slice angka yang awalnya terurut menaik, tapi kemudian dirotasi di suatu titik yang tidak diketahui (misalnya [0,1,2,4,5,6,7] dirotasi jadi [4,5,6,7,0,1,2]). Cari index dari sebuah target di slice tersebut.

// func Search(nums []int, target int) int

// Contoh:
// Input:  nums = [4,5,6,7,0,1,2], target = 0
// Output: 4

// Input:  nums = [4,5,6,7,0,1,2], target = 3
// Output: -1   // tidak ditemukan

// Requirement performa: harus lebih baik dari O(N) — ini yang membuat soal ini menantang, karena array sudah "diacak" (dirotasi), tapi tetap harus dicari lebih cepat dari linear scan.

func Search(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == target {
			return mid
		} else if nums[left] <= nums[mid] {
			if target >= nums[left] && target <= nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			if target >= nums[mid] && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}

	return -1
}

func TestBinarySearch(t *testing.T) {
	nums := []int{4, 5, 6, 7, 0, 1, 2}
	target := 1

	fmt.Println(Search(nums, target))
}
