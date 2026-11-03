package main

// Soal 33: Rentang Log dengan Timestamp Sama

// Ini soal medium, lanjutan dari soal 11 dan 13, dengan satu tantangan baru: angka yang sama bisa muncul berkali-kali.

// Diberikan slice nums yang sudah terurut menaik (misalnya timestamp log, dan beberapa log bisa punya timestamp yang sama) dan sebuah target. Cari index pertama dan index terakhir tempat target muncul. Kalau target tidak ada, return [-1, -1].

// go
// func SearchRange(nums []int, target int) []int

// Contoh:

// Input:  nums = [5,7,7,8,8,10], target = 8
// Output: [3,4]

// Input:  nums = [5,7,7,8,8,10], target = 6
// Output: [-1,-1]

// Input:  nums = [], target = 0
// Output: [-1,-1]

// Requirement performa: harus O(log N). Scan dari kiri atau dari kanan itu O(N), jadi tidak memenuhi.

func SearchRange(nums []int, target int) []int {
	first := findFirst(nums, target)

	if first == -1 {
		return []int{-1, -1}
	}

	last := findLast(nums, target)

	return []int{first, last}
}

func findFirst(nums []int, target int) int {
	left := 0
	right := len(nums) - 1
	result := -1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == target {
			result = mid
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return result
}

func findLast(nums []int, target int) int {
	left := 0
	right := len(nums) - 1
	result := -1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == target {
			result = mid
			left = mid + 1
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return result
}
