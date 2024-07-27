package lrucache

import (
	"sync"
	"time"
)

type Item struct {
	Key     int `json:"key"`
	Value   int `json:"value"`
	ExpTime int `json:"exp"`
}

type LRUCache struct {
	Mu       *sync.Mutex
	Capacity int
	Cache    map[int]*Node
	Head     *Node
	Tail     *Node
}

type Node struct {
	Key     int
	Value   int
	ExpTime time.Time
	Next    *Node
	Prev    *Node
}
