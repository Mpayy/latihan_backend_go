package main

import (
	"log"
	"sync"
	"testing"
)

// Soal 30: Riwayat Produk yang Baru Dilihat

// Ini review LRU cache dari soal 3, dengan studi kasus baru yang sangat umum di e-commerce, area yang kamu kuasai dari proyek microservices-mu: daftar “terakhir dilihat”. API-nya berbeda, jadi jangan langsung menyalin solusi lama.

// Buat komponen yang mencatat produk-produk yang terakhir dilihat satu user, dengan jumlah maksimum tertentu. Kalau sudah penuh, produk yang paling lama tidak dilihat dibuang.

// go
// type RecentlyViewed struct {
//     // isi sendiri
// }

// func NewRecentlyViewed(maxItems int) *RecentlyViewed
// func (r *RecentlyViewed) View(productID string)   // catat bahwa produk ini baru saja dilihat
// func (r *RecentlyViewed) Remove(productID string) // hapus dari riwayat (misalnya user klik "hapus")
// func (r *RecentlyViewed) List() []string          // dari yang paling baru dilihat ke yang paling lama

// Contoh:

// go
// rv := NewRecentlyViewed(3)
// rv.View("A")
// rv.View("B")
// rv.View("C")
// rv.List()        // [C B A]

// rv.View("A")     // A dilihat lagi, naik ke paling baru, tidak muncul dobel
// rv.List()        // [A C B]

// rv.View("D")     // penuh, B (paling lama tidak dilihat) dibuang
// rv.List()        // [D A C]

// rv.Remove("A")
// rv.List()        // [D C]

// Requirement performa: View dan Remove harus O(1). List boleh O(k), dengan k = jumlah produk yang tersimpan.

type NodeItem struct {
	item string
	prev *NodeItem
	next *NodeItem
}

type RecentlyViewed struct {
	mu       sync.Mutex
	maxItems int
	items    map[string]*NodeItem
	head     *NodeItem
	tail     *NodeItem
}

func NewRecentlyViewed(maxItems int) *RecentlyViewed {
	if maxItems <= 0 {
		log.Printf("Warning: maxItems (%d) tidak valid, otomatis diset ke 1", maxItems)
		maxItems = 1
	}

	return &RecentlyViewed{
		maxItems: maxItems,
		items:    make(map[string]*NodeItem),
	}
}

func (r *RecentlyViewed) View(productID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.items[productID]
	if exists {
		r.addToHead(item)
		return
	}

	newNodeItem := &NodeItem{
		item: productID,
	}

	if len(r.items) >= r.maxItems {
		delete(r.items, r.tail.item)
		r.removeNode(r.tail)
	}

	r.addToHead(newNodeItem)
	r.items[productID] = newNodeItem
}

func (r *RecentlyViewed) Remove(productID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.items[productID]
	if !exists {
		return
	}

	delete(r.items, productID)
	r.removeNode(item)
}

func (r *RecentlyViewed) List() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]string, 0, len(r.items))
	curr := r.head

	for curr != nil {
		result = append(result, curr.item)
		curr = curr.next
	}

	return result
}

func (r *RecentlyViewed) addToHead(n *NodeItem) {
	if n == nil {
		return
	}

	if r.head == nil && r.tail == nil {
		r.head = n
		r.tail = n
		return
	}

	if n == r.head {
		return
	}

	if n.prev != nil || n.next != nil || n == r.tail {
		r.removeNode(n)
	}

	n.prev = nil
	n.next = r.head

	r.head.prev = n
	r.head = n
}

func (r *RecentlyViewed) removeNode(n *NodeItem) {
	if n == nil {
		return
	}

	if r.head == n && r.tail == n {
		r.head = nil
		r.tail = nil
		return
	}

	if n == r.head {
		r.head = n.next
		r.head.prev = nil
		return
	}

	if n == r.tail {
		r.tail = n.prev
		r.tail.next = nil
		return
	}

	n.prev.next = n.next
	n.next.prev = n.prev
}



func TestRecentlyViewed(t *testing.T) {
	rv := NewRecentlyViewed(3)
	rv.View("A")
	rv.View("B")
	rv.View("C")
	rv.List() // [C B A]

	rv.View("A") // A dilihat lagi, naik ke paling baru, tidak muncul dobel
	rv.List()    // [A C B]

	rv.View("D") // penuh, B (paling lama tidak dilihat) dibuang
	rv.List()    // [D A C]

	rv.Remove("A")
	rv.List() // [D C]
}
