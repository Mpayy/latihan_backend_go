package main

import (
	"fmt"
	"testing"
)

// Soal 45: Semua Jalur Kategori

// Review rekursi dengan struktur baru: pohon. Katalog e-commerce punya kategori bersarang. Kamu ingin daftar semua jalur dari akar sampai tiap kategori paling ujung (daun, yaitu kategori tanpa anak), misalnya untuk menampilkan breadcrumb.

// go
// type Category struct {
//     Name     string
//     Children []*Category
// }

// func LeafPaths(root *Category) []string

// Return semua jalur dari root ke setiap daun, nama dipisah " > ", urut dari anak kiri ke kanan.

// Contoh:

// Semua
// ├─ Elektronik
// │   ├─ Laptop
// │   │   ├─ Gaming
// │   │   └─ Kantor
// │   └─ HP
// └─ Fashion

// → ["Semua > Elektronik > Laptop > Gaming",
//    "Semua > Elektronik > Laptop > Kantor",
//    "Semua > Elektronik > HP",
//    "Semua > Fashion"]

// Aturan:

// Pohon tidak boleh diubah.
// Hal yang tidak disebut (root nil, elemen nil di Children, root yang tidak punya anak) sengaja tidak dijelaskan. Itu bagian clarifying question-mu.
// Requirement: waktu sebanding dengan ukuran total hasil. Sebutkan apa yang kamu hitung sebagai "ukuran" saat menulis Big-O.

// Format tetap satu pesan: clarifying question, pendekatan, kode, dan Big-O. Di pendekatanmu, jelaskan bagaimana kamu membangun jalur saat turun, dan apa yang terjadi pada jalur itu saat kembali naik ke saudara berikutnya. Kalau rekursi pada pohon terasa asing, tulis "hint 1".

type Category struct {
	Name     string
	Children []*Category
}

func LeafPaths(root *Category) []string {
	if root == nil {
		return nil
	}

	var result []string
	collectLeafPaths(root, "", &result)
	return result
}

func collectLeafPaths(curr *Category, prefix string, result *[]string) {
	if curr == nil {
		return
	}

	currentPath := curr.Name
	if prefix != "" {
		currentPath = prefix + " > " + curr.Name
	}

	if len(curr.Children) == 0 {
		*result = append(*result, currentPath)
		return
	}

	for _, child := range curr.Children {
		collectLeafPaths(child, currentPath, result)
	}
}

func TestLeafPaths(t *testing.T) {
	gaming := &Category{
		Name:     "Gaming",
		Children: []*Category{},
	}
	kantor := &Category{
		Name:     "Kantor",
		Children: []*Category{},
	}
	laptop := &Category{
		Name:     "Laptop",
		Children: []*Category{gaming, kantor},
	}
	hp := &Category{
		Name:     "HP",
		Children: []*Category{},
	}
	elektronik := &Category{
		Name:     "Elektronik",
		Children: []*Category{hp, laptop},
	}
	fashion := &Category{
		Name:     "Fashion",
		Children: []*Category{},
	}
	root := &Category{
		Name:     "Semua",
		Children: []*Category{elektronik, fashion},
	}

	result := LeafPaths(root)
	for _, res := range result {
		fmt.Println(res)
	}
}
