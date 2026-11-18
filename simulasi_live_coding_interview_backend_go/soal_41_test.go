package main

// Soal 41: Berapa Hari Sampai Harga Naik

// Review stack dari soal 10 dengan pola baru. Kamu punya harga harian sebuah produk. Untuk setiap hari, hitung berapa hari lagi sampai muncul harga yang lebih tinggi dari hari itu.

// go
// func DaysUntilHigher(prices []int) []int

// Kalau tidak ada hari yang lebih tinggi setelahnya, isinya 0.

// Contoh:

// [73,74,75,71,69,72,76,73] → [1,1,4,2,1,1,0,0]
// [30,40,50,60]             → [1,1,1,0]
// [30,20,10]                → [0,0,0]

// Requirement: O(N). Brute force dua loop itu O(N²). Format seperti biasa: satu pesan berisi clarifying question, pendekatan, kode, dan Big-O. Di pendekatanmu, sebutkan apa yang disimpan stack (nilai atau index) dan kapan sebuah item dikeluarkan dari stack, dengan alasannya.

func DaysUntilHigher(prices []int) []int {
	result := make([]int, len(prices))
	var stack []int
	for i, v := range prices {
		for len(stack) > 0 && prices[stack[len(stack)-1]] < v {
			prevIndex := stack[len(stack)-1]
			result[prevIndex] = i - prevIndex
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}
	return result
}
