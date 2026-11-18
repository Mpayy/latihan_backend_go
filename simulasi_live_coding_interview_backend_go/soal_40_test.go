package main

import (
	"cmp"
	"container/heap"
	"log"
	"slices"
	"sync"
)

// Soal 40: Pelacak Request Paling Lambat

// Review heap (soal 8 dan 18), dengan kasus yang sangat umum di backend: dari jutaan request yang lewat, kamu hanya ingin menyimpan K request paling lambat, dengan memori yang tidak tumbuh seiring jumlah request.

// go
// type Entry struct {
//     Path      string
//     LatencyMs int
// }

// type SlowTracker struct {
//     // isi sendiri
// }

// func NewSlowTracker(k int) *SlowTracker
// func (t *SlowTracker) Add(path string, latencyMs int)
// func (t *SlowTracker) Top() []Entry // dari paling lambat ke paling cepat

// Aturan:

// Memori O(K), tidak boleh menyimpan semua request yang pernah masuk.
// Add harus O(log K). Top boleh O(K log K).
// Dipakai banyak goroutine sekaligus (thread-safe).
// Latensi sama: urutan bebas.
// Keputusanmu: k <= 0, sebutkan alasannya.
// Top mengembalikan salinan. Perubahan setelahnya tidak boleh memengaruhi tracker, dan sebaliknya.

// Contoh (k = 2): Add ("/a", 50), ("/b", 200), ("/c", 120), ("/d", 10) → Top() = [{/b 200} {/c 120}]. Lalu Add ("/e", 300) → Top() = [{/e 300} {/b 200}].

type Entry struct {
	Path      string
	LatencyMs int
}

type EntryHeap []*Entry

func (h EntryHeap) Len() int           { return len(h) }
func (h EntryHeap) Less(i, j int) bool { return h[i].LatencyMs < h[j].LatencyMs }
func (h EntryHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *EntryHeap) Push(x any) {
	*h = append(*h, x.(*Entry))
}

func (h *EntryHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type SlowTracker struct {
	mu        sync.Mutex
	entryHeap *EntryHeap
	capacity  int
}

func NewSlowTracker(k int) *SlowTracker {
	if k <= 0 {
		log.Printf("[WARN] SlowTracker initialized with invalid k=%d, defaulting to k=1", k)
		k = 1
	}

	slowTracker := &SlowTracker{
		entryHeap: &EntryHeap{},
		capacity:  k,
	}
	heap.Init(slowTracker.entryHeap)
	return slowTracker
}

func (t *SlowTracker) Add(path string, latencyMs int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(*t.entryHeap) < t.capacity {
		heap.Push(t.entryHeap, &Entry{
			Path:      path,
			LatencyMs: latencyMs,
		})
		return
	}

	if (*t.entryHeap)[0].LatencyMs < latencyMs {
		(*t.entryHeap)[0] = &Entry{
			Path:      path,
			LatencyMs: latencyMs,
		}
		heap.Fix(t.entryHeap, 0)
	}
}

func (t *SlowTracker) Top() []Entry {
	t.mu.Lock()
	result := make([]Entry, 0, t.capacity)
	for _, entryHeap := range *t.entryHeap {
		result = append(result, *entryHeap)
	}
	t.mu.Unlock()

	slices.SortFunc(result, func(a, b Entry) int {
		return cmp.Compare(b.LatencyMs, a.LatencyMs)
	})

	return result
}
