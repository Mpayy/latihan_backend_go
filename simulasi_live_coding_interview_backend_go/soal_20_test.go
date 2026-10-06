package main

import (
	"errors"
	"sync"
	"time"
)

// Soal 20 — Simple Circuit Breaker (alternatif library jika mau buat seperti ini adalah gobreaker)

// Soal sistem-oriented lagi — pattern yang sangat umum di backend untuk melindungi sistem dari cascading failure (misalnya saat memanggil service lain yang sedang down, supaya tidak terus-menerus mencoba dan memperburuk keadaan).

// Buat komponen Circuit Breaker sederhana dengan 3 state:

// Closed (normal) — semua request diteruskan ke operation.
// Open (terbuka/terputus) — setelah terlalu banyak kegagalan beruntun, request langsung ditolak tanpa memanggil operation sama sekali, selama periode waktu tertentu.
// Half-Open (coba lagi) — setelah periode "Open" berakhir, satu request diizinkan lewat untuk "test" apakah service sudah sembuh. Kalau berhasil → balik ke Closed. Kalau gagal → balik ke Open lagi (reset timer).
// go
// type CircuitBreaker struct {
//     // isi sendiri
// }

// func NewCircuitBreaker(failureThreshold int, openDuration time.Duration) *CircuitBreaker
// func (cb *CircuitBreaker) Call(operation func() error) error
// failureThreshold: jumlah kegagalan beruntun (consecutive) yang memicu state berubah dari Closed → Open.
// openDuration: berapa lama state Open bertahan sebelum berubah ke Half-Open.
// Call: menjalankan operation (kalau state mengizinkan), return error dari operation, atau error khusus kalau request ditolak karena state Open.

type CircuitBreaker struct {
	mu               sync.Mutex
	state            string
	counter          int
	failureThreshold int
	openDuration     time.Duration
	timestamp        time.Time
}

func NewCircuitBreaker(failureThreshold int, openDuration time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            "Closed",
		failureThreshold: failureThreshold,
		openDuration:     openDuration,
	}
}

var ErrCircuitOpen = errors.New("circuit breaker is open")

func (cb *CircuitBreaker) Call(operation func() error) error {
	cb.mu.Lock()
	initialState := cb.state
	switch initialState {
	case "Closed":
		cb.mu.Unlock()
	case "Open":
		if time.Since(cb.timestamp) < cb.openDuration {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
		initialState = "Half-Open"
		cb.state = "Half-Open"
		cb.mu.Unlock()
	case "Half-Open":
		cb.mu.Unlock()
		return ErrCircuitOpen
	}

	err := operation()

	cb.mu.Lock()
	if err != nil {
		switch initialState {
		case "Half-Open":
			cb.state = "Open"
			cb.timestamp = time.Now()
		case "Closed":
			if cb.state == "Closed" {
				cb.counter++
				if cb.counter >= cb.failureThreshold {
					cb.state = "Open"
					cb.timestamp = time.Now()
				}
			}
		}
		cb.mu.Unlock()
		return err
	}

	switch initialState {
	case "Half-Open":
		cb.state = "Closed"
		cb.counter = 0
	case "Closed":
		cb.counter = 0
	}
	cb.mu.Unlock()
	return nil
}
