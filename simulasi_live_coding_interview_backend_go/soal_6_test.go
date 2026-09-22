package main

import (
	"slices"
	"sync"
	"testing"
)

// Soal 6 — Worker Pool
// "Worker Pool untuk Memproses Jobs"

// Kamu diminta membuat sebuah sistem sederhana yang memproses sejumlah jobs (misalnya: setiap job adalah angka yang harus di-square-kan) menggunakan worker pool — yaitu sejumlah goroutine (worker) tetap yang mengambil job dari sebuah antrian dan memprosesnya secara paralel, alih-alih membuat satu goroutine baru untuk setiap job (yang bisa boros kalau job-nya jutaan).

// func ProcessJobs(jobs []int, numWorkers int) []int

// jobs: daftar angka yang harus diproses (misalnya di-square-kan: x -> x*x).
// numWorkers: jumlah worker (goroutine) yang bekerja secara paralel.
// Return: slice hasil, urutannya harus sama dengan urutan input jobs (meskipun diproses secara paralel dan tidak berurutan).
type toChannel struct {
	Index int
	Value int
}

func ProcessJobs(jobs []int, numWorkers int) []int {
	result := make([]int, len(jobs))
	var wg sync.WaitGroup
	ch := make(chan toChannel)

	go func() {
		for i, v := range jobs {
			ch <- toChannel{
				Index: i,
				Value: v,
			}
		}
		close(ch)
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			for jobs := range ch {
				result[jobs.Index] = jobs.Value * jobs.Value
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return result
}

func TestProcessJobs(t *testing.T) {
	// Definisi kasus uji (Table-Driven Tests)
	tests := []struct {
		name       string
		jobs       []int
		numWorkers int
		expected   []int
	}{
		{
			name:       "Kasus standar dengan 3 worker",
			jobs:       []int{1, 2, 3, 4, 5},
			numWorkers: 3,
			expected:   []int{1, 4, 9, 16, 25},
		},
		{
			name:       "Kasus angka negatif dan nol",
			jobs:       []int{-2, 0, 3},
			numWorkers: 2,
			expected:   []int{4, 0, 9},
		},
		{
			name:       "Jumlah worker lebih banyak dari jumlah jobs",
			jobs:       []int{2, 4},
			numWorkers: 10,
			expected:   []int{4, 16},
		},
		{
			name:       "Worker hanya 1 (sekuensial)",
			jobs:       []int{1, 2, 3},
			numWorkers: 1,
			expected:   []int{1, 4, 9},
		},
		{
			name:       "Jobs kosong",
			jobs:       []int{},
			numWorkers: 4,
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ProcessJobs(tt.jobs, tt.numWorkers)

			if !slices.Equal(res, tt.expected) {
				t.Errorf("Hasil salah untuk %s!\nEkspektasi: %v\nDapatnya   : %v", tt.name, tt.expected, res)
			}
		})
	}
}

// Test khusus untuk memastikan urutan indeks tetap terjaga di bawah tekanan konkurensi tinggi
func TestProcessJobs_OrderPreservation(t *testing.T) {
	n := 1000
	jobs := make([]int, n)
	expected := make([]int, n)

	for i := 0; i < n; i++ {
		jobs[i] = i
		expected[i] = i * i
	}

	numWorkers := 20
	res := ProcessJobs(jobs, numWorkers)

	if !slices.Equal(res, expected) {
		t.Errorf("Urutan indeks tidak sesuai saat diproses oleh %d worker secara paralel!", numWorkers)
	}
}
