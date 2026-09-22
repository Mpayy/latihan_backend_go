package main

import (
	"sync"
	"testing"
	"time"
)

// Soal 1
// "Rate Limiter Sederhana"
// Kamu diminta mengimplementasikan sebuah rate limiter untuk sebuah API service kecil di Go. Rate limiter ini akan dipanggil setiap kali ada request masuk dari seorang user, dan harus memutuskan apakah request tersebut diizinkan atau ditolak.

// Buat sebuah struct/tipe (misalnya RateLimiter) dengan minimal method berikut:

// func NewRateLimiter(limit int, window time.Duration) *RateLimiter
// func (r *RateLimiter) Allow(userID string) bool

// limit: jumlah maksimum request yang diizinkan dalam satu window.
// window: durasi waktu (misalnya 1 menit).
// Allow(userID): dipanggil setiap kali user melakukan request. Return true kalau request diizinkan, false kalau ditolak (melebihi limit).

// Ini harus bisa dipakai oleh banyak user berbeda secara bersamaan (concurrent).

type userBucket struct {
	mu         sync.Mutex
	timestamps []time.Time
}

type RateLimiter struct {
	limit   int
	window  time.Duration
	clients sync.Map
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limit:  limit,
		window: window,
	}
	go rl.CleanUp(window)
	return rl
}

func (r *RateLimiter) CleanUp(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		cutoff := now.Add(-r.window)

		r.clients.Range(func(key, value any) bool {
			userID := key.(string)
			bucket := value.(*userBucket)

			bucket.mu.Lock()

			if len(bucket.timestamps) == 0 || bucket.timestamps[len(bucket.timestamps)-1].Before(cutoff) {
				r.clients.Delete(userID)
			}

			bucket.mu.Unlock()

			return true
		})
	}
}

func (r *RateLimiter) Allow(userID string) bool {
	now := time.Now()
	cutoff := now.Add(-r.window)

	var bucket *userBucket
	if val, exists := r.clients.Load(userID); exists {
		bucket = val.(*userBucket)
	} else {
		newBucket := new(userBucket)
		actual, _ := r.clients.LoadOrStore(userID, newBucket)
		bucket = actual.(*userBucket)
	}

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	validTimestamps := make([]time.Time, 0, len(bucket.timestamps))
	for _, t := range bucket.timestamps {
		if t.After(cutoff) {
			validTimestamps = append(validTimestamps, t)
		}
	}

	if len(validTimestamps) >= r.limit {
		return false
	}

	validTimestamps = append(validTimestamps, now)
	bucket.timestamps = validTimestamps
	return true
}

func TestRateLimiter_AllowLimit(t *testing.T) {
	// Limit 3 request per 1 detik
	rl := NewRateLimiter(3, 1*time.Second)
	userID := "user-1"

	// 3 Request pertama harus diizinkan (true)
	for i := 1; i <= 3; i++ {
		if !rl.Allow(userID) {
			t.Errorf("Request ke-%d harusnya diizinkan (true), tapi malah ditolak (false)", i)
		}
	}

	// Request ke-4 harus ditolak (false) karena sudah capai limit 3
	if rl.Allow(userID) {
		t.Errorf("Request ke-4 harusnya ditolak (false) karena sudah melebihi limit")
	}
}

// 2. Pengujian Reset Window Waktu
func TestRateLimiter_WindowExpiration(t *testing.T) {
	// Limit 2 request per 100 milidetik
	window := 100 * time.Millisecond
	rl := NewRateLimiter(2, window)
	userID := "user-2"

	// Habiskan kuota limit
	rl.Allow(userID)
	rl.Allow(userID)

	if rl.Allow(userID) {
		t.Error("Request ke-3 harusnya ditolak karena limit habis")
	}

	// Tunggu sampai window waktu habis (100ms + buffer 10ms)
	time.Sleep(window + 10*time.Millisecond)

	// Setelah window habis, request harus diizinkan kembali
	if !rl.Allow(userID) {
		t.Error("Request setelah window expired harusnya diizinkan kembali (true)")
	}
}

// 3. Pengujian CleanUp Goroutine
func TestRateLimiter_CleanUp(t *testing.T) {
	// Window 50ms, interval cleanup 50ms
	window := 50 * time.Millisecond
	rl := NewRateLimiter(1, window)
	userID := "user-passive"

	// Buat request agar user tercatat di sync.Map
	rl.Allow(userID)

	// Pastikan data user ada di map
	if _, exists := rl.clients.Load(userID); !exists {
		t.Fatal("User harusnya tersimpan di clients map")
	}

	// Tunggu agar timestamp user basi + CleanUp berjalan (window + interval + buffer)
	time.Sleep(150 * time.Millisecond)

	// Cek apakah CleanUp berhasil menghapus key user dari sync.Map
	if _, exists := rl.clients.Load(userID); exists {
		t.Error("User yang pasif harusnya sudah dihapus oleh CleanUp goroutine")
	}
}

// 4. Pengujian Konkurensi (Aman dari Data Race)
func TestRateLimiter_Concurrency(t *testing.T) {
	limit := 100
	rl := NewRateLimiter(limit, 1*time.Second)
	userID := "user-concurrent"

	var wg sync.WaitGroup
	totalGoroutines := 150
	allowedCount := 0
	var mu sync.Mutex

	// Jalankan 150 goroutine secara bersamaan menembak userID yang sama
	for i := 0; i < totalGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if rl.Allow(userID) {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Hanya boleh ada tepat 100 request yang lolos sesuai `limit`
	if allowedCount != limit {
		t.Errorf("Jumlah request yang diizinkan salah! Ekspektasi %d, tapi dapatnya %d", limit, allowedCount)
	}
}
