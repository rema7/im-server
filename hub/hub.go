package hub

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients	map[*Client]bool
	register  chan *Client
	unregister  chan *Client
	broadcast  chan []byte

}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		clients: make(map[*Client]bool),
		register:   make(chan *Client),
		unregister:   make(chan *Client),
	}
}

func (hub *Hub) Register(conn * websocket.Conn) {
	client := &Client{id: uuid.NewV4().String(), conn: conn, send: make(chan []byte), hub: hub}
	go client.Read()
	go client.Write()

	hub.register <- client
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <- hub.register:
			hub.clients[client] = true
			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			hub.Send(jsonMessage, client)
		case conn := <-hub.unregister:
			if _, ok := hub.clients[conn]; ok {
				close(conn.send)
				delete(hub.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				hub.Send(jsonMessage, conn)
			}

		case message := <-hub.broadcast:
			for conn := range hub.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(hub.clients, conn)
				}
			}
		}
	}
}

func (hub *Hub) Send(message []byte, ignore *Client) {
	for conn := range hub.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}