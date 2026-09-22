package main

import (
	"fmt"
	"testing"
)

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

// 1. Tes Inisialisasi Kuantitas & Error Handling
func TestNewLRUCache(t *testing.T) {
	// Kapasitas valid
	cache, err := NewLRUCache(2)
	if err != nil || cache == nil {
		t.Fatalf("Gagal membuat LRUCache valid: %v", err)
	}

	// Kapasitas invalid (<= 0)
	_, errInvalid := NewLRUCache(0)
	if errInvalid == nil {
		t.Error("Ekspektasi error saat capacity <= 0, tapi tidak ada error")
	}
}

// 2. Tes Operasi Dasar Get & Put
func TestLRUCache_BasicPutGet(t *testing.T) {
	cache, _ := NewLRUCache(2)

	// Put data pertama
	cache.Put(1, 10)

	// Get key yang ada
	val, ok := cache.Get(1)
	if !ok || val != 10 {
		t.Errorf("Ekspektasi Get(1) = (10, true), tapi dapatnya (%d, %t)", val, ok)
	}

	// Get key yang tidak ada
	valNotFound, okNotFound := cache.Get(99)
	if okNotFound || valNotFound != 0 {
		t.Errorf("Ekspektasi Get(99) = (0, false), tapi dapatnya (%d, %t)", valNotFound, okNotFound)
	}
}

// 3. Tes Eviction (Pengosongan Elemen LRU)
func TestLRUCache_Eviction(t *testing.T) {
	cache, _ := NewLRUCache(2)

	cache.Put(1, 100) // Head: 1, Tail: 1
	cache.Put(2, 200) // Head: 2, Tail: 1

	// Mengakses key 1 agar key 1 menjadi Most Recently Used (MRU)
	// Urutan posisi sekarang: Head: 1, Tail: 2 (2 jadi LRU)
	cache.Get(1)

	// Put key 3 (Kapasitas penuh 2/2).
	// Key 2 harus di-evict karena paling jarang dipakai (LRU di Tail)
	cache.Put(3, 300)

	// Pastikan Key 2 sudah terhapus
	if _, ok := cache.Get(2); ok {
		t.Error("Key 2 harusnya sudah di-evict (terhapus), tapi masih ditemukan")
	}

	// Pastikan Key 1 dan Key 3 masih ada
	if val, ok := cache.Get(1); !ok || val != 100 {
		t.Errorf("Key 1 harusnya masih ada, ekspektasi 100 tapi dapat %d", val)
	}

	if val, ok := cache.Get(3); !ok || val != 300 {
		t.Errorf("Key 3 harusnya ada, ekspektasi 300 tapi dapat %d", val)
	}
}

// 4. Tes Update Existing Key
func TestLRUCache_UpdateValue(t *testing.T) {
	cache, _ := NewLRUCache(2)

	cache.Put(1, 100)
	cache.Put(1, 999) // Update value key 1

	val, ok := cache.Get(1)
	if !ok || val != 999 {
		t.Errorf("Value key 1 harusnya terupdate menjadi 999, tapi dapatnya %d", val)
	}

	// Memastikan len(Items) tidak bertambah saat update
	if len(cache.Items) != 1 {
		t.Errorf("Jumlah item di map harusnya tetap 1, tapi dapatnya %d", len(cache.Items))
	}
}
