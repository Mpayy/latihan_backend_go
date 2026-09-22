package main

import "testing"

// Soal 9 — Two Sum
// Ini soal klasik hashmap/array — kemungkinan besar akan muncul dalam bentuk apa pun di interview manapun, jadi worth dipastikan benar-benar solid.

// Diberikan sebuah slice angka nums dan sebuah target. Cari dua index berbeda i dan j sedemikian rupa sehingga nums[i] + nums[j] == target.

// func TwoSum(nums []int, target int) []int

// Contoh:
// Input:  nums = [2,7,11,15], target = 9
// Output: [0,1]   // karena nums[0] + nums[1] = 2 + 7 = 9

func TwoSum(nums []int, target int) []int {
	lastSeen := make(map[int]int)

	for i, num := range nums {
		x := target - num
		if val, ok := lastSeen[x]; ok {
			return []int{val, i}
		}

		lastSeen[num] = i
	}

	return nil
}

func TestTwoSum(t *testing.T) {
	// Inisialisasi data uji (input array dan target jumlah)
	nums := []int{2, 7, 11, 15}
	target := 9

	// Ekspektasi hasil: indeks 0 (angka 2) + indeks 1 (angka 7) = 9
	expected := []int{0, 1}

	// Eksekusi fungsi TwoSum dengan data uji
	res := TwoSum(nums, target)

	// Validasi/Assertion hasil keluaran:
	// Memastikan panjang slice tepat 2 elemen dan nilainya sesuai dengan ekspektasi indeks
	if len(res) != 2 || res[0] != expected[0] || res[1] != expected[1] {
		t.Errorf("Hasil salah! Ekspektasi %v, tapi dapatnya %v", expected, res)
	}
}
