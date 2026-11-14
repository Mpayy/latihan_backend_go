package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

// Soal 36: Hasil Sukses Pertama

// Sebuah halaman produk bisa mengambil harga dari beberapa sumber (cache, replika database, layanan eksternal). Ambil dari semuanya sekaligus, dan pakai yang berhasil paling dulu.

// go
// func FirstSuccess(ctx context.Context, sources []func(ctx context.Context) (string, error)) (string, error)

// Aturan:

// Jalankan semua sources bersamaan.
// Return hasil sukses (err == nil) pertama yang datang, tanpa menunggu yang lain. Sumber lain dibatalkan lewat ctx yang diteruskan.
// Kalau semua gagal, return error. Kamu yang memutuskan bentuknya (error terakhir, error pertama, atau gabungan), sebutkan alasannya.
// Kalau ctx pemanggil dibatalkan sebelum ada yang sukses, return ctx.Err().
// sources kosong: keputusanmu, sebutkan alasannya.
// Setelah FirstSuccess return, goroutine sumber yang masih berjalan boleh selesai sendiri, tapi tidak boleh macet selamanya menunggu penerima yang sudah pergi.

// Contoh: tiga sumber, A sukses "A" setelah 300ms, B gagal setelah 100ms, C sukses "C" setelah 200ms. Hasilnya "C" sekitar 200ms.

type response struct {
	Res string
	Err error
}

func FirstSuccess(ctx context.Context, sources []func(ctx context.Context) (string, error)) (string, error) {
	if len(sources) == 0 {
		return "", errors.New("sources must be greater than 0")
	}

	ctxWC, cancel := context.WithCancel(ctx)
	defer cancel()

	resCh := make(chan response, len(sources))
	var errSlice []error

	for _, source := range sources {
		go func(s func(ctx context.Context) (string, error)) {

			defer func() {
				if r := recover(); r != nil {
					resCh <- response{
						Res: "",
						Err: fmt.Errorf("panic recovered: %v", r),
					}
				}
			}()

			res, err := s(ctxWC)
			resCh <- response{
				Res: res,
				Err: err,
			}
		}(source)
	}

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case res := <-resCh:
			if res.Err == nil {
				return res.Res, nil
			}
			errSlice = append(errSlice, res.Err)
			if len(errSlice) == len(sources) {
				if err := ctx.Err(); err != nil {
					return "", err
				}
				return "", errors.Join(errSlice...)
			}
		}
	}
}

func A(ctx context.Context) (string, error) {
	time.Sleep(300 * time.Millisecond)
	return "A", nil
}

func B(ctx context.Context) (string, error) {
	time.Sleep(100 * time.Millisecond)
	return "", errors.New("source B gagal")
}

func C(ctx context.Context) (string, error) {
	time.Sleep(200 * time.Millisecond)
	return "C", nil
}

func TestFirstSuccess(t *testing.T) {
	sources := []func(ctx context.Context) (string, error){A, B, C}

	res, err := FirstSuccess(t.Context(), sources)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(res)
}
