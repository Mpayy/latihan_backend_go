package main

import (
	"log"
	"math"
	"sync"
	"time"
)

// Soal 31: Rate Limiter Token Bucket

// Ini review rate limiter dari soal 1, dengan algoritma yang dulu sengaja kamu hindari. Waktu itu kamu memilih sliding window karena token bucket terasa "lebih panjang dan rawan salah". Sekarang kita coba, sebagai variasi yang sangat sering muncul sebagai pertanyaan lanjutan di interview.

// Buat rate limiter per user dengan model **token bucket**. Tiap user punya "ember" berisi token:

// - Ember punya kapasitas maksimum `capacity`.
// - Setiap request memakai **1 token**. Kalau ember kosong, request ditolak.
// - Token terisi kembali secara bertahap dengan laju `refillPerSecond` token per detik, tidak pernah melebihi `capacity`.
// - User yang baru pertama kali muncul mulai dengan ember **penuh**.

// ```go
// type TokenBucketLimiter struct {
//     // isi sendiri
// }

// func NewTokenBucketLimiter(capacity int, refillPerSecond float64) *TokenBucketLimiter
// func (l *TokenBucketLimiter) Allow(userID string) bool
// ```

// **Contoh** (`capacity = 3`, `refillPerSecond = 1`):
// ```
// t=0s    Allow("budi") → true   (sisa 2)
// t=0s    Allow("budi") → true   (sisa 1)
// t=0s    Allow("budi") → true   (sisa 0)
// t=0s    Allow("budi") → false  (kosong)
// t=2s    Allow("budi") → true   (sudah terisi 2 token, sisa 1)
// t=2s    Allow("budi") → true   (sisa 0)
// t=2s    Allow("budi") → false
// t=10s   Allow("budi") → true   (terisi penuh lagi, tapi tidak lebih dari 3)
// ```

// **Constraint:** limiter ini dipakai oleh ribuan user sekaligus dan jalan terus berhari-hari. Request dari user yang berbeda tidak boleh saling menunggu, dan **tidak boleh ada goroutine khusus per user** untuk mengisi token.

type userTokenBucket struct {
	mu         sync.Mutex
	tokens     float64
	lastRefill time.Time
}

type TokenBucketLimiter struct {
	mu              sync.RWMutex
	capacity        int
	refillPerSecond float64
	users           map[string]*userTokenBucket
	stopChan        chan struct{}
	once            sync.Once
}

func NewTokenBucketLimiter(capacity int, refillPerSecond float64) *TokenBucketLimiter {
	if capacity <= 0 {
		log.Printf("Warning: capacity (%d) tidak valid, otomatis diset ke 1", capacity)
		capacity = 1
	}

	if refillPerSecond < 0 {
		log.Printf("Warning: refillPerSecond (%f) tidak valid, otomatis diset ke 0", refillPerSecond)
		refillPerSecond = 0
	}

	tb := &TokenBucketLimiter{
		capacity:        capacity,
		refillPerSecond: refillPerSecond,
		users:           make(map[string]*userTokenBucket),
		stopChan:        make(chan struct{}),
	}

	var cleanupInterval time.Duration
	if refillPerSecond == 0 {
		cleanupInterval = 1 * time.Hour
	} else {
		fullRefillSeconds := float64(capacity) / refillPerSecond
		maxIdleDuration := fullRefillSeconds * 1.5

		cleanupInterval = time.Duration(maxIdleDuration * 0.5 * float64(time.Second))

		if cleanupInterval < 1*time.Minute {
			cleanupInterval = 1 * time.Minute
		}
	}

	go tb.StartCleanUp(cleanupInterval)
	return tb
}

func (l *TokenBucketLimiter) Allow(userID string) bool {
	l.mu.RLock()
	bucket, exists := l.users[userID]
	l.mu.RUnlock()

	if !exists {
		l.mu.Lock()
		bucket, exists = l.users[userID]
		if !exists {
			bucket = &userTokenBucket{
				tokens:     float64(l.capacity) - 1,
				lastRefill: time.Now(),
			}
			l.users[userID] = bucket
			l.mu.Unlock()
			return true
		}
		l.mu.Unlock()
	}

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	delta := now.Sub(bucket.lastRefill)
	tokensToAdd := delta.Seconds() * l.refillPerSecond

	bucket.tokens = math.Min(bucket.tokens+tokensToAdd, float64(l.capacity))
	bucket.lastRefill = now

	if bucket.tokens < 1 {
		return false
	}

	bucket.tokens -= 1
	return true
}

func (l *TokenBucketLimiter) StartCleanUp(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.cleanup()
		case <-l.stopChan:
			return
		}
	}
}

func (l *TokenBucketLimiter) cleanup() {
	now := time.Now()

	var maxIdleDuration time.Duration
	if l.refillPerSecond == 0 {
		maxIdleDuration = 1 * time.Hour
	} else {
		fullRefillSeconds := float64(l.capacity) / l.refillPerSecond
		maxIdleDuration = time.Duration(fullRefillSeconds * 1.5 * float64(time.Second))
	}

	l.mu.RLock()
	var expiredKeys []string
	for user, bucket := range l.users {
		bucket.mu.Lock()
		delta := now.Sub(bucket.lastRefill)
		bucket.mu.Unlock()

		if delta > maxIdleDuration {
			expiredKeys = append(expiredKeys, user)
		}
	}
	l.mu.RUnlock()

	if len(expiredKeys) == 0 {
		return
	}

	l.mu.Lock()
	for _, key := range expiredKeys {
		if bucket, exists := l.users[key]; exists {
			bucket.mu.Lock()
			delta := now.Sub(bucket.lastRefill)
			bucket.mu.Unlock()

			if delta > maxIdleDuration {
				delete(l.users, key)
			}

		}
	}
	l.mu.Unlock()
}

func (l *TokenBucketLimiter) Close() {
	l.once.Do(func() {
		close(l.stopChan)
	})
}
