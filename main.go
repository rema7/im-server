package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel



type Hub struct {
	clients	map[*Client]bool
	register  chan *Client
	unregister  chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		register:   make(chan *Client),
		unregister:   make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
			case client := <- h.register:
				h.clients[client] = true
			case client := <- h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
				}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("new connection")
	client := &Client{hub: hub, conn: conn}
	client.hub.register <- client
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	hub := newHub()
	go hub.run()

	http.HandleFunc("/ws",  func(w http.ResponseWriter, r *http.Request) {
		handleConnections(hub, w, r)
	})

	go handleMessages()

	log.Println("http server started on :8100")

	err := http.ListenAndServe(":8100", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
