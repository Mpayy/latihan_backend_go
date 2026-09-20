package main

import (
	"cmp"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

// Soal 1
type userBucket struct {
	mu         sync.Mutex
	timestamps []time.Time
}

type RateLimiter struct {
	limit   int
	window  time.Duration
	clients sync.Map
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limit:  limit,
		window: window,
	}
	go rl.CleanUp(window)
	return rl
}

func (r *RateLimiter) CleanUp(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		cutoff := now.Add(-r.window)

		r.clients.Range(func(key, value any) bool {
			userID := key.(string)
			bucket := value.(*userBucket)

			bucket.mu.Lock()

			if len(bucket.timestamps) == 0 || bucket.timestamps[len(bucket.timestamps)-1].Before(cutoff) {
				r.clients.Delete(userID)
			}

			bucket.mu.Unlock()

			return true
		})
	}
}

func (r *RateLimiter) Allow(userID string) bool {
	now := time.Now()
	cutoff := now.Add(-r.window)

	var bucket *userBucket
	if val, exists := r.clients.Load(userID); exists {
		bucket = val.(*userBucket)
	} else {
		newBucket := new(userBucket)
		actual, _ := r.clients.LoadOrStore(userID, newBucket)
		bucket = actual.(*userBucket)
	}

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	validTimestamps := make([]time.Time, 0, len(bucket.timestamps))
	for _, t := range bucket.timestamps {
		if t.After(cutoff) {
			validTimestamps = append(validTimestamps, t)
		}
	}

	if len(validTimestamps) >= r.limit {
		return false
	}

	validTimestamps = append(validTimestamps, now)
	bucket.timestamps = validTimestamps
	return true
}

// Soal 2
// "Group Anagrams"

// Diberikan sebuah slice of string, kelompokkan string-string yang merupakan anagram satu sama lain ke dalam grup yang sama.
// Dua string dianggap anagram kalau tersusun dari huruf yang sama, hanya urutannya berbeda (misalnya "eat" dan "tea" adalah anagram).

// Signature function:

// func GroupAnagrams(strs []string) [][]string

// Contoh:

// Input:  []string{"eat", "tea", "tan", "ate", "nat", "bat"}
// Output: [][]string{
//     {"eat", "tea", "ate"},
//     {"tan", "nat"},
//     {"bat"},
// }

// urutan grup dan urutan dalam grup tidak masalah

func GroupAnagrams(strs []string) [][]string {
	mapStr := make(map[string][]string, 0)

	for _, str := range strs {
		r := []rune(str)
		slices.Sort(r)
		sortStr := string(r)
		mapStr[sortStr] = append(mapStr[sortStr], str)
	}

	fmt.Println(mapStr)

	sliceStr := make([][]string, 0, len(mapStr))
	for _, vKey := range mapStr {
		sliceStr = append(sliceStr, vKey)
	}

	return sliceStr
}

// Soal 3

// "LRU Cache (Least Recently Used)"
// Implementasikan sebuah cache dengan kapasitas terbatas. Ketika cache penuh dan ada data baru masuk, data yang paling lama tidak diakses (least recently used) harus dibuang otomatis untuk memberi tempat pada data baru.
// Signature yang harus diimplementasikan:

// type LRUCache struct {
//     // isi sendiri
// }

// func NewLRUCache(capacity int) *LRUCache
// func (c *LRUCache) Get(key int) (int, bool)  // return value, dan apakah key ditemukan
// func (c *LRUCache) Put(key int, value int)   // insert atau update

// Requirement performa: Baik Get maupun Put harus berjalan dalam O(1) — bukan O(n). Ini bagian yang membuat soal ini sedikit lebih menantang dari soal sebelumnya.
// Contoh perilaku:

// cache := NewLRUCache(2)
// cache.Put(1, 100)
// cache.Put(2, 200)
// cache.Get(1)        // return 100, true  -> key 1 baru saja diakses, jadi "paling segar"
// cache.Put(3, 300)   // cache penuh (kapasitas 2), maka key 2 dibuang (paling lama tidak diakses)
// cache.Get(2)        // return 0, false   -> key 2 sudah dibuang
// cache.Get(3)        // return 300, true

type Node struct {
	Key   int
	Value int
	Prev  *Node
	Next  *Node
}

type LRUCache struct {
	Capacity int
	Items    map[int]*Node
	Head     *Node
	Tail     *Node
}

func NewLRUCache(capacity int) (*LRUCache, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("capacity must greater than 0")
	}

	return &LRUCache{
		Capacity: capacity,
		Items:    make(map[int]*Node),
	}, nil
}

func (c *LRUCache) Get(key int) (int, bool) {
	val, exists := c.Items[key]
	if !exists || val == nil {
		return 0, false
	}

	c.addToHead(val)

	return val.Value, true
}

func (c *LRUCache) Put(key int, value int) {
	val, exists := c.Items[key]
	if exists {
		val.Value = value
		c.addToHead(val)
		return
	}

	newNode := &Node{
		Key:   key,
		Value: value,
	}

	if len(c.Items) >= c.Capacity {
		delete(c.Items, c.Tail.Key)
		c.removeNode(c.Tail)
	}

	c.Items[key] = newNode
	c.addToHead(newNode)
}

func (c *LRUCache) addToHead(n *Node) {
	if n == nil {
		return
	}

	if c.Head == nil {
		c.Head = n
		c.Tail = n
		return
	}

	if c.Head == n {
		return
	}

	if n.Prev != nil || n.Next != nil || n == c.Tail {
		c.removeNode(n)
	}

	n.Prev = nil
	n.Next = c.Head

	c.Head.Prev = n
	c.Head = n
}

func (c *LRUCache) removeNode(n *Node) {
	if n == nil {
		return
	}

	if n == c.Head && n == c.Tail {
		c.Head = nil
		c.Tail = nil
		return
	}

	if n == c.Head {
		c.Head = n.Next
		c.Head.Prev = nil
		return
	}

	if n == c.Tail {
		c.Tail = n.Prev
		c.Tail.Next = nil
		return
	}

	n.Prev.Next = n.Next
	n.Next.Prev = n.Prev
}

// Soal 4 (Review konsep Soal 1, studi kasus baru)
// "Concurrent Visit Counter per URL"

// Kamu diminta membuat komponen untuk menghitung jumlah visit (kunjungan) ke berbagai URL di sebuah web analytics service. Setiap kali ada request masuk ke suatu URL, counter untuk URL tersebut harus bertambah satu.

// type VisitCounter struct {
//     // isi sendiri
// }

// func NewVisitCounter() *VisitCounter
// func (v *VisitCounter) Increment(url string)
// func (v *VisitCounter) Get(url string) int
// Increment(url): menambah hitungan visit untuk url tersebut sebanyak 1.
// Get(url): mengembalikan jumlah visit saat ini untuk url tersebut (0 kalau belum pernah di-increment).

// Constraint yang harus dipenuhi:
// Service ini akan menerima traffic sangat tinggi dari banyak URL berbeda secara bersamaan (concurrent). Desainmu tidak boleh membuat request ke URL A menunggu request ke URL B — keduanya harus bisa diproses secara paralel tanpa saling blocking, sekalipun banyak goroutine meng-increment URL yang sama di waktu bersamaan juga harus tetap aman (thread-safe, tidak boleh ada data race atau hasil hitungan yang salah).

type VisitCounter struct {
	url sync.Map
}

func NewVisitCounter() *VisitCounter {
	return &VisitCounter{}
}

func (v *VisitCounter) Increment(url string) {
	if val, ok := v.url.Load(url); ok {
		val.(*atomic.Int64).Add(1)
		return
	}

	newCounter := new(atomic.Int64)
	count, _ := v.url.LoadOrStore(url, newCounter)

	countInt64 := count.(*atomic.Int64)
	countInt64.Add(1)
}

func (v *VisitCounter) Get(url string) int {
	count, exists := v.url.Load(url)
	if !exists {
		return 0
	}

	counter := count.(*atomic.Int64)
	return int(counter.Load())
}

// Soal 5
// "Longest Substring Without Repeating Characters"

// Diberikan sebuah string s, cari panjang substring terpanjang yang tidak mengandung karakter berulang.

// func LengthOfLongestSubstring(s string) int

// Contoh:
// Input:  "abcabcbb"
// Output: 3   // substring-nya "abc"

// Input:  "bbbbb"
// Output: 1   // substring-nya "b"

// Input:  "pwwkew"
// Output: 3   // substring-nya "wke"

func LengthOfLongestSubstring(s string) int {
	lastSeen := make(map[rune]int)

	maxLength := 0
	left := 0

	for right, char := range s {
		if lastPos, exists := lastSeen[char]; exists && lastPos >= left {
			left = lastPos + 1
		}

		lastSeen[char] = right

		maxLength = max(maxLength, right-left+1)
	}
	return maxLength
}

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
			fmt.Printf("send to channel: index: %d, value:%d \n", i, v)
			ch <- toChannel{
				Index: i,
				Value: v,
			}
		}
		close(ch)
		fmt.Println("close channel")
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

// Soal 7 — Merge Intervals
// Ini soal array/pattern klasik yang sering muncul juga di kasus nyata backend (misalnya: menggabungkan slot waktu meeting yang bertabrakan, menggabungkan range IP, dll).

// Diberikan sejumlah interval (rentang), masing-masing direpresentasikan sebagai [start, end]. Gabungkan semua interval yang saling tumpang tindih (overlapping) menjadi satu interval saja.

// func MergeIntervals(intervals [][]int) [][]int

// Contoh:
// Input:  [[1,3],[2,6],[8,10],[15,18]]
// Output: [[1,6],[8,10],[15,18]]
// [1,3] dan [2,6] overlap (karena 2 <= 3), digabung jadi [1,6]
// [8,10] dan [15,18] tidak overlap dengan yang lain, tetap terpisah

// Input:  [[1,4],[4,5]]
// Output: [[1,5]]
// [1,4] dan [4,5] dianggap overlap karena bersentuhan di angka 4

func MergeIntervals(intervals [][]int) [][]int {
	var result [][]int

	newIntervals := make([][]int, len(intervals))
	for i, v := range intervals {
		newIntervals[i] = slices.Clone(v)
	}

	slices.SortFunc(newIntervals, func(a, b []int) int {
		if n := cmp.Compare(a[0], b[0]); n != 0 {
			return n
		}

		return cmp.Compare(a[1], b[1])
	})

	for _, v := range newIntervals {
		// result kosong
		if len(result) == 0 {
			result = append(result, v)
			continue
		}

		// overlap
		lastIdx := len(result) - 1
		if v[0] <= result[lastIdx][1] {
			result[lastIdx][1] = max(result[lastIdx][1], v[1])
			continue
		}

		// not overlap
		result = append(result, v)
	}

	return result
}

func main() {
	input := [][]int{{1, 3}, {2, 6}, {8, 10}}
	output := MergeIntervals(input)
	fmt.Println(input) // <- coba tebak, isinya apa SETELAH function ini selesai jalan?
	fmt.Println(output)
}
