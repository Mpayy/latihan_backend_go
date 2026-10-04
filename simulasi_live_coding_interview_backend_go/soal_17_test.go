package main

import (
	"context"
	"testing"
	"time"
)

// Soal 17 — Context Cancellation (Review konsep goroutine, pattern baru)

// Konsep yang perlu dipahami dulu kalau belum familiar: context.Context di Go — ini mekanisme standar untuk membatalkan pekerjaan yang sedang berjalan (misalnya karena timeout, atau karena caller berubah pikiran). Familiar dengan context.Context sebelumnya?

// Kalau belum, saya jelaskan dulu primer singkatnya sebelum soal. Kalau sudah pernah dengar/pakai, langsung saja ke soal:

// Soal: Buat sebuah function yang menjalankan suatu "pekerjaan" (disimulasikan dengan time.Sleep) di goroutine terpisah, tapi berhenti lebih awal kalau context yang diberikan di-cancel (baik karena timeout atau dibatalkan manual), tanpa menunggu pekerjaan itu selesai dulu.

// go
// func DoWorkWithTimeout(ctx context.Context, workDuration time.Duration) error
// Kalau pekerjaan selesai sebelum context di-cancel → return nil.
// Kalau context di-cancel/timeout sebelum pekerjaan selesai → return ctx.Err() (error dari context itu).

func DoWorkWithTimeout(ctx context.Context, workDuration time.Duration) error {
	done := make(chan struct{}, 1)

	go func() {
		time.Sleep(workDuration)
		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func TestDoWorkWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := DoWorkWithTimeout(ctx, 10*time.Second)
	if err != nil {
		t.Log(err)
	}
}
