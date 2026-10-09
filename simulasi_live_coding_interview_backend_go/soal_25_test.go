package main

// Soal 25: Minimum Size Subarray Sum

// Diberikan slice nums berisi angka-angka positif dan sebuah target. Cari panjang subarray berurutan (contiguous) terpendek yang jumlah elemennya lebih besar atau sama dengan target. Kalau tidak ada subarray seperti itu, return 0.

// go
// func MinSubArrayLen(target int, nums []int) int

// Contoh:

// Input:  target = 7, nums = [2,3,1,2,4,3]
// Output: 2
// subarray [4,3] jumlahnya 7, dan itu yang terpendek

// Input:  target = 4, nums = [1,4,4]
// Output: 1
// subarray [4] sudah cukup

// Input:  target = 11, nums = [1,1,1,1,1,1,1,1]
// Output: 0
// jumlah semuanya hanya 8, tidak pernah mencapai 11

func MinSubArrayLen(target int, nums []int) int {
	left := 0
	sum := 0
	minLength := len(nums) + 1

	for right := range len(nums) {
		sum += nums[right]
		for sum >= target {
			currLength := right - left + 1
			if currLength < minLength {
				minLength = currLength
			}
			sum -= nums[left]
			left++
		}
	}

	if minLength == len(nums)+1 {
		minLength = 0
	}

	return minLength
}
