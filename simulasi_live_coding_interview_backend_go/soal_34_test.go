package main

import (
	"sync"
)

// Soal 34: Gabungkan Request yang Sama

// Soal terakhir sesi ini adalah review konkurensi dengan studi kasus baru, dan masih dekat dengan dunia e-commerce yang kamu kenal. Bayangkan halaman produk yang sedang ramai: 100 pengunjung membuka produk yang sama di detik yang sama, dan tanpa penanganan khusus, ke-100-nya memicu 100 query database yang identik.

// Buat komponen yang menggabungkan panggilan bersamaan untuk key yang sama menjadi satu eksekusi:

// go
// type Group struct {
//     // isi sendiri
// }

// func NewGroup() *Group
// func (g *Group) Do(key string, fn func() (string, error)) (string, error)

// Aturan:

// Kalau beberapa goroutine memanggil Do dengan key yang sama sementara fn untuk key itu masih berjalan, fn hanya dijalankan satu kali. Semua pemanggil menunggu lalu menerima hasil yang sama (nilai dan error).
// Setelah fn selesai, key itu dilupakan. Panggilan berikutnya setelah itu menjalankan fn lagi dari awal. Ini bukan cache.
// Key yang berbeda berjalan terpisah dan tidak boleh saling menunggu.

// Contoh:

// go
// g := NewGroup()
// var panggilanDB int32 // dihitung dengan atomic

// loadProduk := func() (string, error) {
//     atomic.AddInt32(&panggilanDB, 1)
//     time.Sleep(100 * time.Millisecond) // simulasi query lambat
//     return "produk-42", nil
// }

// 100 goroutine memanggil bersamaan:
//   g.Do("produk:42", loadProduk)
//
// Hasil: semua menerima ("produk-42", nil), dan panggilanDB == 1

// Satu panggilan lagi setelah semuanya selesai:
//   g.Do("produk:42", loadProduk)
// Hasil: panggilanDB menjadi 2

type call struct {
	ch  chan struct{}
	val string
	err error
}

type Group struct {
	mu    sync.Mutex
	group map[string]*call
}

func NewGroup() *Group {
	return &Group{
		group: make(map[string]*call),
	}
}

func (g *Group) Do(key string, fn func() (string, error)) (string, error) {
	g.mu.Lock()
	c, exists := g.group[key]
	if exists {
		g.mu.Unlock()
		<-c.ch
		return c.val, c.err
	}

	c = &call{
		ch: make(chan struct{}),
	}
	g.group[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()

	g.mu.Lock()
	delete(g.group, key)
	g.mu.Unlock()

	close(c.ch)

	return c.val, c.err
}
