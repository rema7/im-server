package main

import (
	"im-server/cache"
	manager "im-server/hub"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnections(hub *manager.Hub, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query()["token"]

	if len(token) == 0 {
		http.Error(w, "Not authorized", 401)
		return
	}

	redis := cache.GetCache()
	id, err := redis.GetUserId(token[0])
	if err != nil {
		log.Println(err)
		http.Error(w, "Not authorized", 401)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal error", 500)
		return
	}

	log.Println("new connection")
	hub.Register(conn, id)
}

func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	hub := manager.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(hub, w, r)
	})

	log.Println("http server started on :8100")

	err := http.ListenAndServe(":8100", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
