package main

import (
	"container/heap"
	"sync"
	"testing"
)

// Soal 8 — Design Simple In-Memory Job Queue with Priority
// Kali ini gabungan konsep: struct data + sedikit ordering logic, level masih junior/entry tapi lebih ke arah "sistem kecil" daripada algoritma murni.

// Kamu diminta membuat komponen antrian job sederhana, di mana setiap job punya prioritas (angka lebih kecil = lebih prioritas / harus diproses lebih dulu). Ketika job diambil untuk diproses, job dengan prioritas tertinggi (angka terkecil) yang harus keluar duluan.

// type Job struct {
//     ID       string
//     Priority int
// }

// type JobQueue struct {
//     // isi sendiri
// }

// func NewJobQueue() *JobQueue
// func (q *JobQueue) Push(job Job)
// func (q *JobQueue) Pop() (Job, bool)  // ambil & hapus job dengan priority tertinggi (angka terkecil), bool = false kalau queue kosong

type Job struct {
	ID       string
	Priority int
}

type JobHeap []*Job

func (h JobHeap) Len() int           { return len(h) }
func (h JobHeap) Less(i, j int) bool { return h[i].Priority < h[j].Priority }
func (h JobHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *JobHeap) Push(x any) {
	*h = append(*h, x.(*Job))
}

func (h *JobHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type JobQueue struct {
	mu sync.Mutex
	jh *JobHeap
}

func NewJobQueue() *JobQueue {
	jq := &JobQueue{
		jh: &JobHeap{},
	}
	heap.Init(jq.jh)
	return jq
}

func (q *JobQueue) Push(job Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	heap.Push(q.jh, &job)
}

func (q *JobQueue) Pop() (Job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.jh.Len() == 0 {
		return Job{}, false
	}

	res := heap.Pop(q.jh)
	job := res.(*Job)
	return *job, true
}

// 1. Tes Operasi Dasar Push & Urutan Priority (Min-Heap)
func TestJobQueue_Order(t *testing.T) {
	jq := NewJobQueue()

	// Push beberapa job dengan prioritas acak
	jq.Push(Job{ID: "job-medium", Priority: 5})
	jq.Push(Job{ID: "job-low", Priority: 10})
	jq.Push(Job{ID: "job-high", Priority: 1})

	// Ekspektasi urutan Pop: Priority 1 -> 5 -> 10
	expectedOrder := []struct {
		id       string
		priority int
	}{
		{id: "job-high", priority: 1},
		{id: "job-medium", priority: 5},
		{id: "job-low", priority: 10},
	}

	for i, expected := range expectedOrder {
		job, ok := jq.Pop()
		if !ok {
			t.Fatalf("Pop ke-%d gagal! Ekspektasi ada job, tapi dapat false", i+1)
		}

		if job.ID != expected.id || job.Priority != expected.priority {
			t.Errorf("Pop ke-%d salah! Ekspektasi ID=%s Priority=%d, tapi dapat ID=%s Priority=%d",
				i+1, expected.id, expected.priority, job.ID, job.Priority)
		}
	}
}

// 2. Tes Pop pada Queue Kosong
func TestJobQueue_PopEmpty(t *testing.T) {
	jq := NewJobQueue()

	// Pop langsung pada queue kosong
	job, ok := jq.Pop()
	if ok || job != (Job{}) {
		t.Errorf("Ekspektasi Pop() pada queue kosong mengembalikan (Job{}, false), tapi dapat (%v, %t)", job, ok)
	}
}

// 3. Tes Konkurensi (Thread-Safety)
func TestJobQueue_Concurrent(t *testing.T) {
	jq := NewJobQueue()
	totalWorkers := 50
	jobsPerWorker := 100
	totalJobs := totalWorkers * jobsPerWorker

	var wg sync.WaitGroup

	// Push secara simultan dari banyak goroutine
	for i := 0; i < totalWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < jobsPerWorker; j++ {
				jq.Push(Job{
					ID:       "job",
					Priority: j,
				})
			}
		}(i)
	}

	wg.Wait()

	// Lakukan Pop seluruh elemen secara bersamaan dari banyak goroutine
	poppedCount := 0
	var mu sync.Mutex

	for i := 0; i < totalWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				_, ok := jq.Pop()
				if !ok {
					break
				}
				mu.Lock()
				poppedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if poppedCount != totalJobs {
		t.Errorf("Jumlah job yang berhasil di-Pop tidak sesuai! Ekspektasi %d, tapi dapat %d", totalJobs, poppedCount)
	}
}
