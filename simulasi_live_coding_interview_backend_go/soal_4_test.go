package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Soal 4 (Review konsep Soal 1, studi kasus baru)
// "Concurrent Visit Counter per URL"

// Kamu diminta membuat komponen untuk menghitung jumlah visit (kunjungan) ke berbagai URL di sebuah web analytics service. Setiap kali ada request masuk ke suatu URL, counter untuk URL tersebut harus bertambah satu.

// type VisitCounter struct {
//     // isi sendiri
// }

// func NewVisitCounter() *VisitCounter
// func (v *VisitCounter) Increment(url string)
// func (v *VisitCounter) Get(url string) int
// Increment(url): menambah hitungan visit untuk url tersebut sebanyak 1.
// Get(url): mengembalikan jumlah visit saat ini untuk url tersebut (0 kalau belum pernah di-increment).

// Constraint yang harus dipenuhi:
// Service ini akan menerima traffic sangat tinggi dari banyak URL berbeda secara bersamaan (concurrent). Desainmu tidak boleh membuat request ke URL A menunggu request ke URL B — keduanya harus bisa diproses secara paralel tanpa saling blocking, sekalipun banyak goroutine meng-increment URL yang sama di waktu bersamaan juga harus tetap aman (thread-safe, tidak boleh ada data race atau hasil hitungan yang salah).

type VisitCounter struct {
	url sync.Map
}

func NewVisitCounter() *VisitCounter {
	return &VisitCounter{}
}

func (v *VisitCounter) Increment(url string) {
	if val, ok := v.url.Load(url); ok {
		val.(*atomic.Int64).Add(1)
		return
	}

	newCounter := new(atomic.Int64)
	count, _ := v.url.LoadOrStore(url, newCounter)

	countInt64 := count.(*atomic.Int64)
	countInt64.Add(1)
}

func (v *VisitCounter) Get(url string) int {
	count, exists := v.url.Load(url)
	if !exists {
		return 0
	}

	counter := count.(*atomic.Int64)
	return int(counter.Load())
}

// 1. Tes Operasi Dasar Increment & Get
func TestVisitCounter_Basic(t *testing.T) {
	vc := NewVisitCounter()
	url := "https://example.com"

	// Cek URL yang belum pernah di-increment (harus 0)
	if count := vc.Get("https://notfound.com"); count != 0 {
		t.Errorf("Ekspektasi count = 0 untuk URL tidak dikenal, tapi dapatnya %d", count)
	}

	// Increment 3 kali
	vc.Increment(url)
	vc.Increment(url)
	vc.Increment(url)

	expected := 3
	res := vc.Get(url)

	if res != expected {
		t.Errorf("Hasil salah! Ekspektasi count = %d, tapi dapatnya %d", expected, res)
	}
}

// 2. Tes Multiple URL
func TestVisitCounter_MultipleURLs(t *testing.T) {
	vc := NewVisitCounter()

	urlA := "https://site.com/a"
	urlB := "https://site.com/b"

	vc.Increment(urlA)
	vc.Increment(urlA)
	vc.Increment(urlB)

	if countA := vc.Get(urlA); countA != 2 {
		t.Errorf("Ekspektasi count URL A = 2, tapi dapatnya %d", countA)
	}

	if countB := vc.Get(urlB); countB != 1 {
		t.Errorf("Ekspektasi count URL B = 1, tapi dapatnya %d", countB)
	}
}

// 3. Tes Konkurensi (Aman dari Data Race & Lost Updates)
func TestVisitCounter_Concurrent(t *testing.T) {
	vc := NewVisitCounter()
	url := "https://concurrent.com"

	totalGoroutines := 100
	incrementsPerGoroutine := 1000
	expectedTotal := totalGoroutines * incrementsPerGoroutine

	var wg sync.WaitGroup

	// Jalankan 100 goroutine secara bersamaan, masing-masing melakukan increment 1000 kali
	for i := 0; i < totalGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				vc.Increment(url)
			}
		}()
	}

	wg.Wait()

	res := vc.Get(url)
	if res != expectedTotal {
		t.Errorf("Hasil konkurensi salah! Ekspektasi total count = %d, tapi dapatnya %d", expectedTotal, res)
	}
}
