package hub

import (
	"github.com/gorilla/websocket"
)

type Result struct {
	id      int
	message []byte
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Result
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Result),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) Register(conn *websocket.Conn) {
	client := &Client{id: -1, conn: conn, send: make(chan []byte), hub: hub}

	go client.Read()
	go client.Write()

	hub.register <- client
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case conn := <-hub.unregister:
			if _, ok := hub.clients[conn]; ok {
				close(conn.send)
				delete(hub.clients, conn)
			}

		case result := <-hub.broadcast:
			for conn := range hub.clients {
				select {
				case conn.send <- result.message:
				default:
					close(conn.send)
					delete(hub.clients, conn)
				}
			}
		}
	}
}

// func (hub *Hub) Send(message []byte, ignore *Client) {
// 	for conn := range hub.clients {
// 		if conn != ignore {
// 			conn.send <- message
// 		}
// 	}
// }
