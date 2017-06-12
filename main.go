package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"encoding/json"
	"github.com/satori/go.uuid"
)

type Hub struct {
	clients	map[*Client]bool
	register  chan *Client
	unregister  chan *Client
	broadcast  chan []byte

}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		clients: make(map[*Client]bool),
		register:   make(chan *Client),
		unregister:   make(chan *Client),
	}
}

func (hub *Hub) run() {
	for {
		select {
			case client := <- hub.register:
				hub.clients[client] = true
				jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
				hub.send(jsonMessage, client)
		case conn := <-hub.unregister:
			if _, ok := hub.clients[conn]; ok {
				close(conn.send)
				delete(hub.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				hub.send(jsonMessage, conn)
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

func (hub *Hub) send(message []byte, ignore *Client) {
	for conn := range hub.clients {
		if conn != ignore {
			conn.send <- message
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
	client := &Client{id: uuid.NewV4().String(), conn: conn, send: make(chan []byte), hub: hub}
	hub.register <- client
	go client.read()
	go client.write()
}



func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	hub := newHub()
	go hub.run()

	http.HandleFunc("/ws",  func(w http.ResponseWriter, r *http.Request) {
		handleConnections(hub, w, r)
	})

	//go handleMessages()

	log.Println("http server started on :8100")

	err := http.ListenAndServe(":8100", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
