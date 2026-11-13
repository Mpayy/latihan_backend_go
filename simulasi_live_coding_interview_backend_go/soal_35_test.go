package main

import "sync"

// Soal 35: Gabungkan Beberapa Aliran Data

// Bayangkan sebuah dashboard yang mengambil data dari beberapa sumber sekaligus (misalnya layanan stok, harga, dan promo), masing-masing menyalurkan hasilnya lewat channel. Kamu perlu menggabungkannya menjadi satu aliran.

// func Merge(inputs ...<-chan int) <-chan int

// Aturan:

// Channel hasil berisi semua nilai dari semua channel di inputs.
// Channel hasil harus ditutup setelah semua input ditutup dan semua nilainya sudah diteruskan.
// Urutan antar channel berbeda tidak ditentukan. Urutan nilai dari channel yang sama harus terjaga.
// Merge harus langsung return channel hasilnya, tidak boleh menunggu.
// Tidak boleh ada goroutine yang tertinggal setelah semuanya selesai.

// Contoh:

// a := make(chan int)
// b := make(chan int)
// go func() { a <- 1; a <- 2; close(a) }()
// go func() { b <- 10; close(b) }()

// for v := range Merge(a, b) {
//     fmt.Println(v)
// }
// mencetak 1, 2, 10 (10 boleh muncul di mana saja, tapi 1 selalu sebelum 2)
// lalu loop berhenti sendiri karena channel hasil ditutup

func Merge(inputs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range inputs {
		wg.Add(1)
		go func(c <-chan int) {
			for hasil := range c {
				out <- hasil
			}
			wg.Done()
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
