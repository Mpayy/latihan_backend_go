package main

import (
	"sync"
)

// Soal 12 — Simple Pub/Sub System
// Kamu diminta membuat komponen Publish-Subscribe (Pub/Sub) sederhana in-memory. Beberapa "subscriber" bisa mendaftar untuk menerima pesan pada suatu "topic", dan ketika ada yang "publish" pesan ke topic tersebut, semua subscriber yang terdaftar di topic itu harus menerima pesan tersebut.

// type PubSub struct {
//     // isi sendiri
// }

// func NewPubSub() *PubSub
// func (ps *PubSub) Subscribe(topic string) <-chan string
// func (ps *PubSub) Publish(topic string, message string)

// Subscribe(topic): mendaftar sebagai subscriber ke topic, return sebuah channel yang akan menerima pesan-pesan yang di-publish ke topic itu setelah ini.
// Publish(topic, message): mengirim message ke SEMUA subscriber yang terdaftar di topic tersebut.

// Constraint: Boleh ada banyak subscriber untuk satu topic yang sama. Publish tidak boleh macet/blocking selamanya hanya karena ada satu subscriber yang lambat/tidak pernah membaca channel-nya.

type PubSub struct {
	mu          sync.Mutex
	subscribers map[string][]chan string
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan string),
	}
}

func (ps *PubSub) Subscribe(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

func (ps *PubSub) Publish(topic string, message string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, ch := range ps.subscribers[topic] {
		select {
		case ch <- message:
		default:
		}
	}
}
