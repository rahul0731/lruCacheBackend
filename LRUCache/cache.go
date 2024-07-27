package lrucache

import (
	"errors"
	"sync"
	"time"
)

func NewLRUCache(capacity int, mu *sync.Mutex) *LRUCache {
	return &LRUCache{
		Capacity: capacity,
		Cache:    make(map[int]*Node, capacity),
		Mu:       mu,
	}
}

func (l *LRUCache) AddToFront(node *Node) {
	node.Next = l.Head
	if l.Head != nil {
		l.Head.Prev = node
	}
	l.Head = node
	if l.Tail == nil {
		l.Tail = node
	}
}

func (l *LRUCache) DeleteNode(node *Node) {
	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		l.Tail = node.Prev
	}
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		l.Head = node.Next
	}
}

func (l *LRUCache) GetItem(key int) (int, error) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	node, exist := l.Cache[key]
	if !exist {
		return -1, errors.New("item not exist")
	}
	if node.ExpTime.Before(time.Now()) {
		l.DeleteNode(node)
		return -1, errors.New("item not exist")

	}
	l.DeleteNode(node)
	l.AddToFront(node)
	return node.Value, nil
}

func (l *LRUCache) SetItem(item Item) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	ExpTime := time.Now().Add(time.Second * time.Duration(item.ExpTime))
	node, exist := l.Cache[item.Key]
	if exist {
		node.Value = item.Value
		node.ExpTime = ExpTime
		l.DeleteNode(node)
		l.AddToFront(node)
		return nil
	}
	newNode := &Node{Key: item.Key, Value: item.Value, ExpTime: ExpTime}
	if len(l.Cache) >= l.Capacity {
		delete(l.Cache, l.Tail.Key)
		l.DeleteNode(l.Tail)
	}
	l.Cache[item.Key] = newNode
	l.AddToFront(newNode)
	return nil
}

func (l *LRUCache) DeleteItem(key int) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	node, exist := l.Cache[key]
	if !exist {
		return errors.New("item not exist")
	}
	delete(l.Cache,key)
	l.DeleteNode(node)
	return nil
}

func (l *LRUCache) ClearExpired() {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	for key, node := range l.Cache {
		if node.ExpTime.Before(time.Now()) {
			delete(l.Cache, key)
			l.DeleteNode(node)
		}
	}
}

func (l *LRUCache) ExpirationChecker(t int) {
	go func() {
		for {
			time.Sleep(time.Duration(t) * time.Second)
			l.ClearExpired()
		}
	}()
}
