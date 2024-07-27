package main

import (
	lrucache "LRU-cache/LRUCache"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan lrucache.Item)
)

func main() {
	var mu sync.Mutex
	cache := lrucache.NewLRUCache(3, &mu)
	cache.ExpirationChecker(60)
	controller := lrucache.NewLRUCntroller(cache, broadcast)

	r := mux.NewRouter()
	r.HandleFunc("/set", controller.SetItemController).Methods("POST")
	r.HandleFunc("/get", controller.GetItemController).Methods("GET")
	r.HandleFunc("/delete", controller.DeleteItemController).Methods("DELETE")
	r.HandleFunc("/ws", handleConnections)

	corsOptions := handlers.AllowedOrigins([]string{"*"})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"})

	go handleMessages()

	fmt.Println("Listening and serving on :8080 ....")
	http.ListenAndServe(":8080", handlers.CORS(corsOptions, corsHeaders, corsMethods)(r))

}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	fmt.Println("New WebSocket connection")

	for {
		var item lrucache.Item
		err := ws.ReadJSON(&item)
		if err != nil {
			fmt.Println("WebSocket read error:", err)
			delete(clients, ws)
			break
		}
		fmt.Printf("Received item: %+v\n", item)
		broadcast <- item
	}
}

func handleMessages() {
	for {
		item := <-broadcast
		fmt.Printf("Broadcasting item: %+v\n", item)
		for client := range clients {
			err := client.WriteJSON(item)
			if err != nil {
				fmt.Println("WebSocket write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
