package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Soal 26: In-Memory Cache dengan TTL

// Review konsep concurrency dan cleanup dari soal 1 (rate limiter), dengan studi kasus berbeda yang sangat umum di backend nyata: cache yang datanya kedaluwarsa otomatis.

// Buat cache in-memory sederhana. Setiap data disimpan dengan TTL (time-to-live), yaitu berapa lama data itu dianggap valid. Setelah TTL lewat, data dianggap tidak ada lagi.

// go
// type TTLCache struct {
//     // isi sendiri
// }

// func NewTTLCache() *TTLCache
// func (c *TTLCache) Set(key string, value string, ttl time.Duration)
// func (c *TTLCache) Get(key string) (string, bool)  // bool = false kalau key tidak ada ATAU sudah kedaluwarsa

// Contoh perilaku:

// go
// cache := NewTTLCache()
// cache.Set("token", "abc123", 100*time.Millisecond)

// cache.Get("token")   // "abc123", true

// // ... 150ms kemudian
// cache.Get("token")   // "", false (sudah kedaluwarsa)

// Constraint: cache ini akan dipakai oleh banyak goroutine sekaligus, dan service-nya jalan terus-menerus berhari-hari.

type cacheItem struct {
	value     string
	expiredAt time.Time
}

type TTLCache struct {
	mu       sync.RWMutex
	cache    map[string]*cacheItem
	stopChan chan struct{}
	once     sync.Once
}

func NewTTLCache() *TTLCache {
	c := &TTLCache{
		cache:    make(map[string]*cacheItem),
		stopChan: make(chan struct{}),
	}
	go c.StartCleanUp(60 * time.Second)
	return c
}

func (c *TTLCache) StartCleanUp(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopChan:
			return
		}
	}
}

func (c *TTLCache) cleanup() {
	c.mu.RLock()
	now := time.Now()
	var expiredKeys []string

	for key, item := range c.cache {
		if now.After(item.expiredAt) {
			expiredKeys = append(expiredKeys, key)
		}
	}
	c.mu.RUnlock()

	if len(expiredKeys) == 0 {
		return
	}

	c.mu.Lock()
	for _, key := range expiredKeys {
		if item, exists := c.cache[key]; exists {
			if now.After(item.expiredAt) {
				delete(c.cache, key)
			}
		}
	}
	c.mu.Unlock()
}

func (c *TTLCache) Close() {
	c.once.Do(func() {
		close(c.stopChan)
	})
}

func (c *TTLCache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &cacheItem{
		value:     value,
		expiredAt: time.Now().Add(ttl),
	}
}

func (c *TTLCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.cache[key]
	if !exists {
		return "", false
	}

	if time.Now().After(item.expiredAt) {
		return "", false
	}

	return item.value, true
}

func TestTTLCache(t *testing.T) {
	cache := NewTTLCache()
	cache.Set("token", "abc123", 100*time.Millisecond)

	valueTrue, getTrue := cache.Get("token") // "abc123", true
	fmt.Println(valueTrue, getTrue)

	time.Sleep(150 * time.Millisecond)
	valueFalse, getFalse := cache.Get("token") // "", false (sudah kedaluwarsa)
	fmt.Println(valueFalse, getFalse)
}
