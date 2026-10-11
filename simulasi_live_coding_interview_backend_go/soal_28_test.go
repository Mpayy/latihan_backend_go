package main

import (
	"fmt"
	"sync"
	"testing"
)

// Soal 28: Process All dengan Error per Item

// Soal penutup sesi ini adalah review worker pool dari soal 6. Soal 6 dulu butuh banyak pancingan, jadi sekarang kita lihat seberapa jauh polanya sudah melekat. Ada variasi baru: sekarang pekerjaan bisa gagal.

// Diberikan slice inputs dan sebuah function fn yang memproses satu input. fn bisa mengembalikan error. Jalankan fn untuk semua input secara paralel memakai sejumlah worker tetap, lalu kumpulkan hasilnya.

// go
// type Result struct {
//     Value int
//     Err   error
// }

// func ProcessAll(inputs []int, numWorkers int, fn func(int) (int, error)) []Result
// Hasilnya harus berurutan sesuai inputs, meski pemrosesannya paralel.
// Kegagalan satu input tidak boleh menghentikan input lain. Setiap input mendapat Result-nya sendiri.

// Contoh:

// go
// inputs := []int{1, 2, 3, 4}
// fn := func(x int) (int, error) {
//     if x%2 == 0 {
//         return 0, fmt.Errorf("angka genap: %d", x)
//     }
//     return x * 10, nil
// }

// ProcessAll(inputs, 2, fn)
// [{10 nil} {0 "angka genap: 2"} {30 nil} {0 "angka genap: 4"}]

type Result struct {
	Value int
	Err   error
}

type toProcess struct {
	index int
	value int
}

func ProcessAll(inputs []int, numWorkers int, fn func(int) (int, error)) []Result {
	if len(inputs) == 0 {
		return []Result{}
	}

	if numWorkers <= 0 {
		numWorkers = 1
	}

	result := make([]Result, len(inputs))
	ch := make(chan toProcess)
	var wg sync.WaitGroup

	go func() {
		for i, v := range inputs {
			ch <- toProcess{
				index: i,
				value: v,
			}
		}
		close(ch)
	}()

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for jobs := range ch {
				res, err := fn(jobs.value)
				result[jobs.index] = Result{
					Value: res,
					Err:   err,
				}
			}
		}()
	}

	wg.Wait()

	return result
}

func TestProcessAll(t *testing.T) {
	inputs := []int{1,2,3,4}
	fn := func(x int) (int, error) {
		if x%2 == 0 {
			return 0, fmt.Errorf("angka genap: %d", x)
		}
		return x * 10, nil
	}

	res := ProcessAll(inputs, -1, fn)
	fmt.Println(res)
}
