package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	manager "im/hub"
)



var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnections(hub *manager.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("new connection")
	hub.Register(conn)
}



func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	hub := manager.NewHub()
	go hub.Run()

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
