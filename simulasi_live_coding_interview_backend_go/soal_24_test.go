package main

import (
	"sync"
	"time"
)

// Soal 24 — Debounce Function

// Soal praktis/sistem-oriented untuk menutup sesi ini — pattern yang umum dipakai di backend (misalnya menunda eksekusi suatu aksi sampai tidak ada "trigger" baru dalam periode waktu tertentu, mirip seperti debounce di frontend tapi versi backend/Go).

// Buat sebuah function yang membungkus sebuah function lain (fn), sehingga setiap kali dipanggil, eksekusi fn yang sebenarnya ditunda selama delay. Kalau function hasil bungkusan itu dipanggil lagi sebelum delay habis, timer sebelumnya dibatalkan dan diulang dari awal — jadi fn hanya benar-benar dieksekusi kalau sudah tidak ada pemanggilan baru selama delay.

// go
// func Debounce(fn func(), delay time.Duration) func()

// Contoh ilustrasi (alur, bukan kode):

// debounced := Debounce(doSomething, 100*time.Millisecond)

// debounced()  // panggilan 1, di t=0ms
// debounced()  // panggilan 2, di t=50ms -> timer dari panggilan 1 dibatalkan
// debounced()  // panggilan 3, di t=90ms -> timer dari panggilan 2 dibatalkan

// doSomething() BARU benar-benar jalan di sekitar t=190ms
// (100ms setelah panggilan TERAKHIR, yaitu panggilan 3)

func Debounce(fn func(), delay time.Duration) func() {
	var timer *time.Timer
	var mu sync.Mutex
	return func() {
		mu.Lock()
		defer mu.Unlock()

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(delay, fn)
	}
}
