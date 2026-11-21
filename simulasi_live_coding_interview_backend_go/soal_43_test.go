package main

import (
	"errors"
	"fmt"
	"testing"
)

// Soal 43: Urutan Eksekusi Job dengan Dependensi

// Konsep baru dan belum pernah kita latih: graf. Sebuah build system atau job scheduler harus menjalankan job dalam urutan yang menghormati dependensi.

// Job adalah node, dan dependensi adalah panah dari job yang harus dijalankan lebih dulu ke job yang bergantung padanya. Dua istilah yang perlu kamu kenal: in-degree (jumlah panah yang masuk ke sebuah node) dan siklus (lingkaran dependensi, misalnya A butuh B dan B butuh A).

// go
// func JobOrder(jobs []string, deps map[string][]string) ([]string, error)

// deps["B"] = ["A"] berarti B bergantung pada A, jadi A harus dijalankan lebih dulu. Return urutan eksekusi yang valid dan memuat semua jobs. Kalau ada siklus, return error.

// Contoh:

// jobs = ["compile","test","package","deploy"]
// deps = {"test":["compile"], "package":["compile","test"], "deploy":["package"]}
// → ["compile","test","package","deploy"]

// jobs = ["a","b"], deps = {"a":["b"], "b":["a"]}
// → error (siklus)

// Aturan:

// Job yang sama-sama siap boleh dalam urutan apa saja, asal hasilnya valid.
// Requirement: O(V + E), dengan V = jumlah job dan E = jumlah dependensi.
// Hal-hal lain yang tidak disebut (misalnya dependensi ke job yang tidak ada di jobs, atau job yang muncul dua kali) sengaja tidak saya jelaskan. Itu bagian clarifying question-mu.

// Format tetap satu pesan: clarifying question, pendekatan, kode, dan Big-O dengan variabel yang didefinisikan. Di pendekatanmu, jelaskan bagaimana kamu tahu sebuah job sudah siap dijalankan dan bagaimana kamu mendeteksi siklus. Kalau konsep graf terasa asing, tulis “hint 1”.

func JobOrder(jobs []string, deps map[string][]string) ([]string, error) {
	inDegree := make(map[string]int, len(jobs))
	dependents := make(map[string][]string, len(jobs))
	queue := make([]string, 0, len(jobs))
	result := make([]string, 0, len(jobs))

	for _, job := range jobs {
		if _, ok := inDegree[job]; ok {
			return nil, errors.New("duplicated job in jobs")
		}
		inDegree[job] = 0
	}

	for job, prereq := range deps {
		if _, ok := inDegree[job]; !ok {
			return nil, errors.New("missing job in deps")
		}

		seen := make(map[string]struct{})
		for _, req := range prereq {
			if _, ok := inDegree[req]; !ok {
				return nil, errors.New("missing dependency in jobs")
			}

			if _, ok := seen[req]; !ok {
				seen[req] = struct{}{}
				inDegree[job]++
				dependents[req] = append(dependents[req], job)
			}
		}
	}

	for _, job := range jobs {
		if inDegree[job] == 0 {
			queue = append(queue, job)
		}
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		result = append(result, curr)
		nextJobs := dependents[curr]
		for _, nextJob := range nextJobs {
			inDegree[nextJob]--
			if inDegree[nextJob] == 0 {
				queue = append(queue, nextJob)
			}
		}
	}

	if len(result) != len(jobs) {
		return nil, errors.New("cycle detected in job dependencies")
	}

	return result, nil
}

func TestJobOrder(t *testing.T) {
	// jobs := []string{"compile", "test", "package", "deploy"}
	// jobs := []string{"deploy", "package", "test", "compile"}
	// jobs := []string{"a", "b"}
	jobs := []string{"a", "a"}
	// deps := map[string][]string{
	// 	"test":    {"compile"},
	// 	"package": {"compile", "test"},
	// 	"deploy":  {"package"},
	// }
	// deps := map[string][]string{
	// 	"a": {"b"},
	// 	"b": {"a"},
	// }
	deps := map[string][]string{
		"z": {"a"},
	}

	res, err := JobOrder(jobs, deps)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}
