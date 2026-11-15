package main

import (
	"context"
	"errors"
	"time"
)

// Soal 39: Pengumpul Batch

// Kamu menerima aliran event (log, klik, pesanan) dan ingin menulisnya ke database per batch, bukan satu per satu, supaya lebih hemat.

// go
// func Batch(ctx context.Context, in <-chan int, maxSize int, maxWait time.Duration, flush func(batch []int)) error

// Aturan:

// Kumpulkan nilai dari in ke dalam batch.
// Panggil flush kalau batch sudah berisi maxSize item, atau maxWait sudah lewat sejak item pertama batch itu masuk, mana yang lebih dulu.
// Setelah flush, mulai batch baru. Batch kosong tidak pernah di-flush.
// in ditutup: flush sisa batch (kalau ada) lalu return nil.
// ctx dibatalkan: return ctx.Err(). Kamu yang memutuskan apakah sisa batch di-flush dulu, sebutkan alasannya.
// flush dipanggil dari goroutine yang menjalankan Batch dan boleh menyimpan slice yang diterimanya. Pikirkan apa artinya itu untuk batch berikutnya.
// Tidak ada goroutine yang tertinggal setelah Batch return.
// Keputusanmu: maxSize <= 0 dan maxWait <= 0.

// Contoh (maxSize = 3, maxWait = 100ms): kirim 1, 2, 3 → flush [1 2 3]. Kirim 4, lalu diam → flush [4] sekitar 100ms setelah 4 masuk. Kirim 5, 6, lalu tutup in → flush [5 6] dan Batch return nil.

func Batch(ctx context.Context, in <-chan int, maxSize int, maxWait time.Duration, flush func(batch []int)) error {
	if maxSize <= 0 {
		return errors.New("maxSize must be greater than 0")
	}

	if maxWait <= 0 {
		return errors.New("maxWait must be greater than 0")
	}

	batch := make([]int, 0, maxSize)
	timer := time.NewTimer(maxWait)
	timer.Stop()

	doFlush := func() {
		if len(batch) > 0 {
			flush(batch)
			batch = make([]int, 0, maxSize)
		}
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}

	for {
		select {
		case event, ok := <-in:
			if !ok {
				doFlush()
				return nil
			}

			if len(batch) == 0 {
				timer.Reset(maxWait)
			}

			batch = append(batch, event)
			if len(batch) == maxSize {
				doFlush()
			}

		case <-timer.C:
			doFlush()

		case <-ctx.Done():
			doFlush()
			return ctx.Err()
		}
	}
}
