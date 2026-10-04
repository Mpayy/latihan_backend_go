package main

import "time"

// Soal 19 — Retry with Exponential Backoff

// Soal praktis/sistem-oriented lagi, sesuai preferensimu — dan ini pattern yang sangat umum dipakai di backend nyata (misalnya saat memanggil API eksternal yang kadang gagal sementara).

// Buat sebuah function yang mencoba ulang (retry) sebuah operasi yang bisa gagal, dengan jeda waktu yang semakin lama di antara setiap percobaan (exponential backoff), sampai berhasil atau mencapai batas maksimum percobaan.

// go
// func RetryWithBackoff(operation func() error, maxRetries int, baseDelay time.Duration) error
// operation: function yang mencoba melakukan sesuatu, return error kalau gagal, nil kalau berhasil.
// maxRetries: jumlah maksimum percobaan ulang (retry) kalau operasi terus gagal.
// baseDelay: jeda waktu dasar sebelum percobaan pertama. Setiap percobaan berikutnya, jeda waktu berlipat ganda (misal: baseDelay, 2*baseDelay, 4*baseDelay, dst).
// Return nil kalau operation berhasil (di percobaan manapun). Return error terakhir kalau sudah mencoba sampai maxRetries kali dan tetap gagal semua.

// Contoh ilustrasi (alur, bukan kode):

// baseDelay = 100ms, maxRetries = 3

// Percobaan 1: operation() gagal → tunggu 100ms
// Percobaan 2: operation() gagal → tunggu 200ms
// Percobaan 3: operation() gagal → tunggu 400ms
// Percobaan 4: operation() berhasil → return nil

func RetryWithBackoff(operation func() error, maxRetries int, baseDelay time.Duration) error {
	var lastErr error
	for i := 0; i <= maxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		if i < maxRetries {
			time.Sleep(baseDelay * (1 << i))
		}
	}

	return lastErr
}
