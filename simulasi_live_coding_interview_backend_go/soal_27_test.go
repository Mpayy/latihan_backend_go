package main

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// Soal 27: Top Error Endpoints (Parsing Log)

// Soal ini tentang parsing data, yang termasuk kategori praktis yang kamu minta di awal dan belum pernah kita latih. Kasusnya sangat umum di backend: membaca log request dan merangkumnya.

// Diberikan sebuah slice berisi baris log. Setiap baris berformat:

// <METHOD> <PATH> <STATUS> <LATENCY_MS>

// Cari k endpoint (path) yang paling sering mengembalikan error server, yaitu status 500 ke atas. Urutkan dari yang paling banyak error. Return daftar path-nya saja.

// go
// func TopErrorEndpoints(logs []string, k int) []string

// Contoh:

// logs = [
//   "GET /api/users 200 35",
//   "GET /api/orders 500 120",
//   "POST /api/orders 503 300",
//   "GET /api/users 404 12",
//   "GET /api/orders 502 90",
//   "GET /api/payments 500 800",
// ]
// k = 2

// Output: ["/api/orders", "/api/payments"]
// /api/orders punya 3 error server, /api/payments punya 1
// /api/users tidak masuk hitungan: 200 dan 404 bukan error server

type pathCount struct {
	path  string
	count int
}

func TopErrorEndpoints(logs []string, k int) []string {
	if k <= 0 {
		return nil
	}

	errorCounts := make(map[string]int)
	for _, line := range logs {
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			code, err := strconv.Atoi(fields[2])
			if err == nil && code >= 500 {
				errorCounts[fields[1]]++
			}
		}
	}

	fullPath := make([]pathCount, 0, len(errorCounts))
	for path, count := range errorCounts {
		fullPath = append(fullPath, pathCount{
			path:  path,
			count: count,
		})
	}

	slices.SortFunc(fullPath, func(a, b pathCount) int {
		return cmp.Compare(b.count, a.count)
	})

	if k > len(fullPath) {
		k = len(fullPath)
	}

	result := make([]string, 0, k)
	for i := 0; i < k; i++ {
		result = append(result, fullPath[i].path)
	}

	return result
}

func TestTopErrorEndpoints(t *testing.T) {
	logs := []string{
		"GET /api/users 200 35",
		"GET /api/orders 500 120",
		"POST /api/orders 503 300",
		"GET /api/users 404 12",
		"GET /api/orders 502 90",
		"GET /api/payments 500 800",
	}

	k := -1

	result := TopErrorEndpoints(logs, k)
	fmt.Println(result)
}
