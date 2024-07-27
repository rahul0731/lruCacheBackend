package lrucache

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type LRUController struct {
	lruCache  *LRUCache
	broadcast chan Item
}

func NewLRUCntroller(lru *LRUCache, broadcast chan Item) *LRUController {
	return &LRUController{lruCache: lru, broadcast: broadcast}
}

func (l *LRUController) GetItemController(w http.ResponseWriter, r *http.Request) {
	keyString := r.URL.Query().Get("key")
	key, err := strconv.Atoi(keyString)
	if err != nil {
		http.Error(w, "Key not converted to int", http.StatusBadRequest)
		return
	}
	value, err := l.lruCache.GetItem(key)

	if err != nil {
		http.Error(w, "Key not found :"+err.Error(), http.StatusNotFound)
		return
	}
	result := fmt.Sprintf("Key : %d, Value : %d", key, value)
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusOK)

}

func (l *LRUController) SetItemController(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := l.lruCache.SetItem(item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := fmt.Sprintf("Item added to cache, Key : %d", item.Key)
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusOK)
	l.broadcast <- item
}

func (l *LRUController) DeleteItemController(w http.ResponseWriter, r *http.Request) {
	keyString := r.URL.Query().Get("key")
	key, err := strconv.Atoi(keyString)
	if err != nil {
		http.Error(w, "Key not converted to int", http.StatusBadRequest)
		return
	}
	err = l.lruCache.DeleteItem(key)
	if err != nil {
		http.Error(w, "Key not found :"+err.Error(), http.StatusNotFound)
		return
	}
	result := fmt.Sprintf("Item deleted, Key : %d", key)
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusOK)

}
