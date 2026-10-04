package main

import (
	"cmp"
	"fmt"
	"slices"
	"testing"
)

// Soal 18 — Kth Largest Element (Review Heap, studi kasus baru)

// Review konsep heap dari soal 8 (Priority Job Queue), tapi dengan kasus yang berbeda — supaya pattern-nya makin melekat, bukan cuma hafal satu soal.

// Diberikan sebuah slice angka nums dan sebuah angka k. Cari elemen terbesar ke-k dalam slice itu (bukan elemen unik ke-k, duplikat dihitung).

// go
// func FindKthLargest(nums []int, k int) int

// Contoh:

// Input:  nums = [3,2,1,5,6,4], k = 2
// Output: 5   // terbesar ke-1 = 6, terbesar ke-2 = 5

// Input:  nums = [3,2,3,1,2,4,5,5,6], k = 4
// Output: 4

func FindKthLargest(nums []int, k int) int {
	slices.SortFunc(nums, func(a, b int) int {
		return cmp.Compare(b, a)
	})

	return nums[k-1]
}

func TestFindKthLargest(t *testing.T) {
	nums := []int{3, 2, 3, 1, 2, 4, 5, 5, 6}
	result := FindKthLargest(nums, 4)
	fmt.Println(result)
}
