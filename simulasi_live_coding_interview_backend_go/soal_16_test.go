package main

import (
	"testing"
	"time"
)

// Soal 16 — Bounded Concurrency (Semaphore Sederhana)

// Soal penutup sesi ini — kembali ke konteks praktis/sistem-oriented, dan sekaligus review konsep channel/goroutine dari soal 6 (worker pool) dengan pattern yang sedikit berbeda.

// Kamu diminta membuat komponen yang membatasi jumlah operasi yang berjalan bersamaan (concurrent), meskipun ada banyak goroutine yang ingin menjalankan operasi itu di waktu yang sama. Misalnya: kamu punya 100 request masuk bersamaan, tapi hanya boleh maksimal 5 yang benar-benar diproses di saat yang sama (sisanya harus menunggu giliran).

// go
// type Semaphore struct {
//     // isi sendiri
// }

// func NewSemaphore(maxConcurrent int) *Semaphore
// func (s *Semaphore) Acquire()  // menunggu sampai ada "slot" kosong, lalu mengambil slot itu
// func (s *Semaphore) Release()  // mengembalikan slot, supaya goroutine lain bisa jalan

// Contoh penggunaan (ilustrasi, tidak perlu ditulis):

// go
// sem := NewSemaphore(5)

// for i := 0; i < 100; i++ {
//     go func() {
//         sem.Acquire()
//         defer sem.Release()
//         // ... lakukan sesuatu, maksimal 5 goroutine yang jalan ke sini bersamaan
//     }()
// }

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(maxConcurrent int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, maxConcurrent),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func TestSemaphore(t *testing.T) {
	sem := NewSemaphore(5)

	for i := 0; i < 100; i++ {
		sem.Acquire()
		go func() {
			defer sem.Release()
			t.Logf("goroutine %d masuk antrian", i)
			time.Sleep(time.Second)
		}()
	}
}
