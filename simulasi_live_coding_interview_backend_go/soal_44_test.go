package main

import (
	"fmt"
	"testing"
)

// Soal 44: Jalur Terpendek di Gudang

// Review graf dengan pola baru. Robot gudang bergerak di grid. Sel 0 adalah lorong, sel 1 adalah rak (penghalang). Robot bergerak satu langkah ke atas, bawah, kiri, atau kanan (tanpa diagonal), hanya lewat lorong.

// go
// func ShortestPath(grid [][]int, start, target [2]int) int

// start dan target berbentuk [baris, kolom]. Return jumlah langkah minimum dari start ke target, atau -1 kalau tidak bisa dicapai.

// Contoh:

// grid:  0 0 0
//        1 1 0
//        0 0 0

// start=[0,0], target=[2,0] → 6   // harus memutar lewat kolom kanan
// start=[0,0], target=[0,0] → 0

// grid:  0 1
//        1 0
// start=[0,0], target=[1,1] → -1  // diagonal tidak boleh

// Requirement: O(R × C), dengan R baris dan C kolom. Hal yang tidak disebut (misalnya start atau target di atas rak, di luar grid, atau grid kosong) sengaja tidak dijelaskan, itu bagian clarifying question-mu.

// Format tetap satu pesan: clarifying question, pendekatan, kode, dan Big-O. Di pendekatanmu, jelaskan kenapa BFS, bukan DFS, menjamin jalur terpendek, dan bagaimana kamu mencegah satu sel dikunjungi dua kali. Kalau BFS terasa asing, tulis "hint 1".

func ShortestPath(grid [][]int, start, target [2]int) int {
	R := len(grid)
	if R == 0 || len(grid[0]) == 0 {
		return -1
	}
	C := len(grid[0])

	inBounds := func(pos [2]int) bool {
		return pos[0] >= 0 && pos[0] < R && pos[1] >= 0 && pos[1] < C
	}

	if !inBounds(start) || !inBounds(target) {
		return -1
	}

	if grid[start[0]][start[1]] == 1 || grid[target[0]][target[1]] == 1 {
		return -1
	}

	if start[0] == target[0] && start[1] == target[1] {
		return 0
	}

	visited := make([][]bool, R)
	for i := range visited {
		visited[i] = make([]bool, C)
	}
	visited[start[0]][start[1]] = true

	type node struct {
		r, c, dist int
	}

	queue := make([]node, 0)
	queue = append(queue, node{
		r:    start[0],
		c:    start[1],
		dist: 0,
	})

	dirs := [][2]int{
		{-1, 0}, // atas
		{1, 0},  // bawah
		{0, -1}, // kiri
		{0, 1},  // kanan
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr.r == target[0] && curr.c == target[1] {
			return curr.dist
		}

		for _, d := range dirs {
			nr := curr.r + d[0]
			nc := curr.c + d[1]

			if !inBounds([2]int{nr, nc}) {
				continue
			}

			if grid[nr][nc] == 1 || visited[nr][nc] {
				continue
			}

			visited[nr][nc] = true
			queue = append(queue, node{r: nr, c: nc, dist: curr.dist + 1})
		}
	}

	return -1
}

func TestShortestPath(t *testing.T) {
	// grid := [][]int{
	// 	{0, 0, 0},
	// 	{1, 1, 0},
	// 	{0, 0, 0},
	// }
	grid := [][]int{
		{0, 1},
		{1, 0},
	}

	// start := [2]int{0, 0}
	// target := [2]int{0, 0}
	start := [2]int{0, 0}
	target := [2]int{1, 1}

	res := ShortestPath(grid, start, target)
	fmt.Println(res)
}
