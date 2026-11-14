package main

import (
	"context"
	"errors"
	"sync"
)

// Soal 37: Paralel dengan Batas, Berhenti di Error Pertama (pola semaphore)

// Mengimpor 1000 produk dari file, tapi database hanya boleh menerima 5 penulisan sekaligus, dan begitu satu gagal, sisanya tidak perlu dilanjutkan.

// go
// func ForEachParallel(ctx context.Context, items []int, limit int, fn func(ctx context.Context, item int) error) error

// Aturan:

// Jalankan fn untuk setiap item, maksimal limit yang berjalan bersamaan.
// Begitu satu fn gagal: batalkan context yang diteruskan ke fn lain, jangan memulai item baru, dan return error pertama yang terjadi.
// Return hanya setelah semua goroutine yang sudah dimulai selesai. Tidak ada yang tertinggal (beda dengan soal 36).
// ctx pemanggil dibatalkan: berhenti memulai item baru, return ctx.Err() kalau belum ada error lain.
// Semua sukses: nil. limit <= 0 dan items kosong: keputusanmu, sebutkan alasannya.

// Contoh: items = [1..6], limit = 2. fn menunggu 100ms (atau berhenti saat ctx dibatalkan), kecuali item 3 yang gagal setelah 50ms. Item 1-2 jalan di 0-100ms, item 3-4 mulai di 100ms, item 3 gagal di 150ms dan item 4 berhenti karena pembatalan. Item 5 dan 6 tidak pernah dijalankan, dan hasilnya error dari item 3.

func ForEachParallel(ctx context.Context, items []int, limit int, fn func(ctx context.Context, item int) error) error {
	if limit <= 0 {
		return errors.New("limit must be greater than 0")
	}

	ctxWC, cancel := context.WithCancel(ctx)
	defer cancel()

	chErr := make(chan error, 1)
	var wg sync.WaitGroup

	sem := make(chan struct{}, limit)

	for _, item := range items {
		if err := ctxWC.Done(); err != nil {
			break
		}

		wg.Add(1)
		go func(job int) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
				defer func() {
					<-sem
				}()
			case <-ctxWC.Done():
				return
			}

			if err := ctxWC.Done(); err != nil {
				return
			}

			if err := fn(ctxWC, job); err != nil {
				select {
				case chErr <- err:
					cancel()
				default:
				}
			}

		}(item)

	}
	wg.Wait()

	select {
	case err := <-chErr:
		return err
	default:
		if err := ctx.Err(); err != nil {
			return err
		}
		return nil
	}
}
